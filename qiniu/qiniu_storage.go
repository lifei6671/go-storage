package qiniu

import (
	"io"
	"github.com/lifei6671/go-storage"
	"github.com/qiniu/api.v7/auth/qbox"
	qstorage "github.com/qiniu/api.v7/storage"
	"errors"
	scontext "context"
	"bytes"
	"fmt"
)

var (
	ErrBucketNoEmpty = errors.New("bucket name does not empty")
	ErrNotSupportedFileAppend = errors.New("file append is not supported")
	ErrNotSupportedCreateDir = errors.New("creating a directory is not supported")
)

type QiniuStorage struct {
	mac *qbox.Mac
	bucketManager *qstorage.BucketManager
}

func NewQiniuStorage(accessKey string,secretKey string,cfg qstorage.Config) *QiniuStorage  {

	mac := qbox.NewMac(accessKey, secretKey)

	bucketManager := qstorage.NewBucketManager(mac, &cfg)

	return &QiniuStorage{
		bucketManager : bucketManager,
		mac : mac,
	}
}

//获取文件信息
func (s *QiniuStorage) Stat(path string,context storage.Context) (storage.FileInfo,error){

	bucket,ok := context.Get("bucket")

	if !ok {
		return nil,ErrBucketNoEmpty
	}
	fileInfo, err := s.bucketManager.Stat(bucket.(string), path)
	if err != nil {
		return nil,err
	}

	info := NewQiniuFileInfo(fileInfo)
	info.name = path

	return info,nil
}

//写入文件
func (s *QiniuStorage) Write(path string, contents string, context storage.Context) (int64,error){
	bucket,ok := context.Get("bucket")

	if !ok {
		return 0,ErrBucketNoEmpty
	}
	putPolicy := qstorage.PutPolicy{
		Scope:               bucket.(string),
	}

	upToken := putPolicy.UploadToken(s.mac)


	cfg := qstorage.Config { UseHTTPS :false,UseCdnDomains:false}

	if c,ok := context.Get("config");ok{
		if cc,ok := c.(qstorage.Config);ok {
			cfg = cc
		}
	}

	// 构建表单上传的对象
	formUploader := qstorage.NewFormUploader(&cfg)

	ret := qstorage.PutRet{}
	// 可选配置
	putExtra := qstorage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	buf := bytes.NewBufferString(contents)

	dataLen := int64(buf.Len())

	err := formUploader.Put(scontext.Background(), &ret, upToken, path,buf,dataLen, &putExtra)
	if err != nil {

		return 0,err
	}
	return dataLen,nil
}

//以流的方式写入到七牛云储存中
func (s *QiniuStorage) WriteStream(path string, reader io.Reader,size int64, context storage.Context) (error){
	bucket,ok := context.Get("bucket")

	if !ok {
		return ErrBucketNoEmpty
	}
	putPolicy := qstorage.PutPolicy{
		Scope:               bucket.(string),
	}

	upToken := putPolicy.UploadToken(s.mac)


	cfg := qstorage.Config { UseHTTPS :false,UseCdnDomains:false}

	if c,ok := context.Get("config");ok{
		if cc,ok := c.(qstorage.Config);ok {
			cfg = cc
		}
	}

	// 构建表单上传的对象
	formUploader := qstorage.NewFormUploader(&cfg)

	ret := qstorage.PutRet{}
	// 可选配置
	putExtra := qstorage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.Put(scontext.Background(), &ret, upToken, path,reader,size, &putExtra)
	if err != nil {

		return err
	}
	return nil
}

//以字节的方式写入到七牛云储存
func (s *QiniuStorage) WriteBytes(path string, contents []byte, context storage.Context) (int64,error){
	bucket,ok := context.Get("bucket")

	if !ok {
		return 0,ErrBucketNoEmpty
	}
	putPolicy := qstorage.PutPolicy{
		Scope:               bucket.(string),
	}

	upToken := putPolicy.UploadToken(s.mac)


	cfg := qstorage.Config { UseHTTPS :false,UseCdnDomains:false}

	if c,ok := context.Get("config");ok{
		if cc,ok := c.(qstorage.Config);ok {
			cfg = cc
		}
	}

	// 构建表单上传的对象
	formUploader := qstorage.NewFormUploader(&cfg)

	ret := qstorage.PutRet{}
	// 可选配置
	putExtra := qstorage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	buf := bytes.NewBuffer(contents)

	dataLen := int64(buf.Len())

	err := formUploader.Put(scontext.Background(), &ret, upToken, path,buf,dataLen, &putExtra)
	if err != nil {

		return 0,err
	}
	return dataLen,nil
}

func (s *QiniuStorage) Append(path string, contents string, context storage.Context) (int64,error){
	return 0, ErrNotSupportedFileAppend
}

func (s *QiniuStorage) AppendStream(path string, reader io.Reader, context storage.Context) (int64,error){
	return 0, ErrNotSupportedFileAppend
}

func (s *QiniuStorage) AppendBytes(path string, contents []byte, context storage.Context) (int64,error){
	return 0, ErrNotSupportedFileAppend
}

//重命名文件
func (s *QiniuStorage) ReName(path string,newPath string,context storage.Context) error{
	srcBucket := ""
	if sb,ok := context.Get("srcBucket"); ok {
		if sbt,ok := sb.(string); ok {
			srcBucket = sbt
		}
	}
	destBucket := ""
	if sb,ok := context.Get("destBucket"); ok {
		if sbt,ok := sb.(string); ok {
			destBucket = sbt
		}
	}
	force := false
	if sb,ok := context.Get("destBucket"); ok {
		if sbt,ok := sb.(bool); ok {
			force = sbt
		}
	}

	err := s.bucketManager.Move(srcBucket, path, destBucket, newPath, force)

	return err

}

//复制文件
func (s *QiniuStorage) Copy (src string,dst string, context storage.Context) error{
	srcBucket := ""
	if sb,ok := context.Get("srcBucket"); ok {
		if sbt,ok := sb.(string); ok {
			srcBucket = sbt
		}
	}
	destBucket := ""
	if sb,ok := context.Get("destBucket"); ok {
		if sbt,ok := sb.(string); ok {
			destBucket = sbt
		}
	}
	force := false
	if sb,ok := context.Get("destBucket"); ok {
		if sbt,ok := sb.(bool); ok {
			force = sbt
		}
	}

	err := s.bucketManager.Copy(srcBucket, src, destBucket, dst, force)

	return err
}

//删除文件
func (s *QiniuStorage) Delete(dst string,context storage.Context) error{
	bucket := ""
	if sb,ok := context.Get("bucket"); ok {
		if sbt,ok := sb.(string); ok {
			bucket = sbt
		}
	}

	err := s.bucketManager.Delete(bucket,dst)

	return err
}

//删除目录，目前实现的是删除指定前缀的文件
func (s *QiniuStorage) DeleteDir(dst string,context storage.Context) error{
	bucket := ""
	if sb,ok := context.Get("bucket"); ok {
		if sbt,ok := sb.(string); ok {
			bucket = sbt
		}
	}
	delimiter := ""
	marker := ""
	limit := 1000


	for {
		entries, _, nextMarker, hashNext, err := s.bucketManager.ListFiles(bucket, dst, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}

		deleteOps := make([]string, 0, len(entries))

		//print entries
		for _, entry := range entries {
			//fmt.Println(entry.Key)
			deleteOps = append(deleteOps, qstorage.URIDelete(bucket, entry.Key))

		}

		rets, err := s.bucketManager.Batch(deleteOps)
		if err != nil {

			// 遇到错误
			if _, ok := err.(*qstorage.ErrorInfo); ok {
				for _, ret := range rets {
					// 200 为成功
					if ret.Code != 200 {
						fmt.Printf("%s\n", ret.Data.Error)
						return errors.New(ret.Data.Error)
					}
				}
			} else {
				fmt.Printf("batch error, %s", err)
			}
		}

		if hashNext {
			marker = nextMarker
		}
	}
	return nil
}

func (s *QiniuStorage) CreateDir(dst string, context storage.Context) error{
	return ErrNotSupportedCreateDir
}

//列出指定前缀的文件
func (s *QiniuStorage) ListDir(dst string,ctx storage.Context) ([]storage.FileInfo,int64,error) {
	bucket := ""
	if sb,ok := ctx.Get("bucket"); ok {
		if sbt,ok := sb.(string); ok {
			bucket = sbt
		}
	}
	delimiter := ""
	marker := ""
	limit := 1000

	files := make([]storage.FileInfo,0)
	num := 0

	for {
		entries, _, nextMarker, hashNext, err := s.bucketManager.ListFiles(bucket, dst, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}
		num = num + len(entries)

		deleteOps := make([]string, 0, len(entries))

		for _, entry := range entries {
			deleteOps = append(deleteOps, qstorage.URIDelete(bucket, entry.Key))

			info := qstorage.FileInfo{
				Hash:entry.Hash,
				Fsize:entry.Fsize,
				PutTime: entry.PutTime,
				MimeType : entry.MimeType,
				Type: entry.Type,
			}

			info.Hash = entry.Hash

			fileInfo := NewQiniuFileInfo(info)
			fileInfo.name = entry.Key

			files = append(files,fileInfo)

		}


		if hashNext {
			marker = nextMarker
		}
	}
	return files ,int64(num) ,nil
}
