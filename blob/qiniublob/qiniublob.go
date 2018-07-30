package qiniublob

import (
	"github.com/lifei6671/go-storage/blob"
	"github.com/qiniu/api.v7/storage"
	"context"
	"github.com/qiniu/api.v7/auth/qbox"
	"errors"
	"github.com/lifei6671/go-storage/blob/driver"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"net/http"
	"time"
	"os"
	"crypto/md5"
	"encoding/hex"
	"bytes"
)

var emptyBody = ioutil.NopCloser(strings.NewReader(""))

var QiniuDeadline = time.Now().Add(time.Second * 3600).Unix() //1小时有效期

type reader struct {
	body        io.ReadCloser
	size        int64
	contentType string
}

func (r *reader) Attrs() *driver.ObjectAttrs {
	return &driver.ObjectAttrs{
		Size:        r.size,
		ContentType: r.contentType,
	}
}

func (r *reader) Read(p []byte) (int, error) {
	return r.body.Read(p)
}

// Close closes the reader itself. It must be called when done reading.
func (r *reader) Close() error {
	return r.body.Close()
}

type bucket struct {
	mac    *qbox.Mac
	cfg    *storage.Config
	scope  string
	domain string
}

type qiniuiError struct {
	bucket, key, msg string
	kind             driver.ErrorKind
}

func (e qiniuiError) Error() string {
	return fmt.Sprintf("gcs://%s/%s: %s", e.bucket, e.key, e.msg)
}
func (e qiniuiError) BlobError() driver.ErrorKind {
	return e.kind
}

func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (driver.Reader, error) {

	if offset < 0 {
		return nil, fmt.Errorf("negative offset %d", offset)
	}

	bucketManager := storage.NewBucketManager(b.mac, b.cfg)

	fileInfo, sErr := bucketManager.Stat(b.scope, key)
	if sErr != nil {
		return nil, qiniuiError{b.scope, key, sErr.Error(), driver.NotFound}
	}

	if length == 0 {

		return &reader{size: fileInfo.Fsize, contentType: fileInfo.MimeType, body: emptyBody}, nil
	}

	privateAccessURL := storage.MakePrivateURL(b.mac, b.domain, key, QiniuDeadline)

	client := http.Client{}

	req, err := http.NewRequest("GET", privateAccessURL, nil)

	if offset > 0 && length < 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	} else if length > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
	}

	if err != nil {
		return nil, qiniuiError{bucket: b.scope, key: key, msg: err.Error(), kind: driver.NotFound}
	}
	resp, err := client.Do(req)

	return &reader{size: fileInfo.Fsize, contentType: fileInfo.MimeType, body: resp.Body}, nil
}

func (b *bucket) NewTypedWriter(ctx context.Context, key string, contentType string, opt *driver.WriterOptions) (driver.Writer, error) {

	putPolicy := storage.PutPolicy{
		Scope: b.scope,
	}
	upToken := putPolicy.UploadToken(b.mac)

	//如果没有缓冲区，则会将每次写入的数据写入到内存，直到调用Close方法
	return &writer{
		toggle:      &blob.Toggle{},
		cfg:         b.cfg,
		ctx:         ctx,
		bufferSize:  opt.BufferSize,
		key:         key,
		size:        0,
		contentType: contentType,
		upToken:     upToken,
		uploader:	 storage.NewResumeUploader(b.cfg),
		isFirst:	 true,
	}, nil

}

func (b *bucket) Delete(ctx context.Context, key string) error {
	bucketManager := storage.NewBucketManager(b.mac, b.cfg)

	return bucketManager.Delete(b.scope,key)
}

func OpenBucket( accessKey string, secretKey string, cfg storage.Config, scope string, domain string) (*blob.Bucket, error) {
	if accessKey == "" {
		return nil, errors.New("accessKey must be provided to get bucket")
	}
	if secretKey == "" {
		return nil, errors.New("secretKey must be provided to get bucket")
	}

	return blob.NewBucket(&bucket{
		scope:  scope,
		mac:    qbox.NewMac(accessKey, secretKey),
		cfg:    &cfg,
		domain: domain,
	}), nil
}

type writer struct {
	toggle      *blob.Toggle
	w 			*io.PipeWriter
	r 			*io.PipeReader
	isFirst		bool
	bucket      string
	f           *os.File
	key         string
	bufferSize  int
	upHost      string
	size        int64
	ctx         context.Context
	cfg         *storage.Config
	contentType string
	uploader    *storage.ResumeUploader
	upToken     string
	ret         *storage.BlkputRet
	err         error
}

func (w *writer) Write(p []byte) (int, error) {
	//如果设置了缓冲，则分片上传
	if w.bufferSize > 0 {
		dataLength := len(p)
		var err error
		w.toggle.Do(func() {
			pr,pw := io.Pipe()
			w.w = pw
			w.r = pr
			w.w.Write(p)

			scheme := "http://"
			if w.cfg.UseHTTPS {
				scheme = "https://"
			}

			host := w.cfg.Zone.SrcUpHosts[0]
			if w.cfg.UseCdnDomains {
				host = w.cfg.Zone.CdnUpHosts[0]
			}
			upHost := fmt.Sprintf("%s%s", scheme, host)
			w.upHost = upHost
			var ret storage.BlkputRet


			err = w.uploader.Mkblk(w.ctx, w.upToken, upHost, &ret, w.bufferSize, bytes.NewReader(p), -1)

		}, func() {
			ret := storage.BlkputRet{}
			w.size += int64(dataLength)
			err = w.uploader.Bput(w.ctx, w.upToken, &ret, bytes.NewReader(p), dataLength)
		})

		return dataLength,err
	} else {
		l := 0
		var err error
		w.toggle.Do(func() {
			f, e := ioutil.TempFile(os.TempDir(), "qiniu")
			if e != nil {
				err = e
			} else {
				w.f = f
			}
		}, func() {
			l,err = w.f.Write(p)
		})

		if w.err != nil {
			return 0, w.err
		}

		return l, nil
	}
}

func (w *writer) Close() error {
	ret := storage.PutRet{}

	if w.f != nil {
		path := w.f.Name()
		fmt.Println(path)

		defer func() {
			if err := os.Remove(path); err != nil {
				w.err = err
			}
		}()

		if err := w.f.Close(); err != nil {
			w.err = err
		}

		putExtra := storage.RputExtra{}
		if err := w.uploader.PutFile(w.ctx,&ret,w.upToken,w.key,path,&putExtra); err != nil {
			return err
		}
		return nil
	}
	if w.bufferSize > 0 {
		w.err = w.uploader.Mkfile(w.ctx, w.upToken, w.upHost, &ret, w.key, false, w.size, nil)
	}
	return w.err
}

func md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
