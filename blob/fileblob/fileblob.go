package fileblob

import (
	"os"
	"fmt"
	"github.com/lifei6671/go-storage/blob"
	"strings"
	"path/filepath"
	"context"
	slashpath "path"
	"io"
	"github.com/lifei6671/go-storage/blob/driver"
)

type bucket struct {
	dir string
}

// NewBucket creates a new bucket that reads and writes to dir.
// dir must exist.
func NewBucket(dir string) (*blob.Bucket, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("open file bucket: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("open file bucket: %s is not a directory", dir)
	}
	return blob.NewBucket(&bucket{dir}), nil
}

// resolvePath converts a key into a relative filesystem path. It guarantees
// that there will only be one valid key for a given path and that the resulting
// path will not reach outside the directory.
func resolvePath(key string) (string, error) {
	for _, c := range key {
		if !('A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '/' || c == '.' || c == ' ' || c == '_' || c == '-') {
			return "", fmt.Errorf("contains invalid character %q", c)
		}
	}
	if cleaned := slashpath.Clean(key); key != cleaned {
		return "", fmt.Errorf("not a clean slash-separated path")
	}
	if slashpath.IsAbs(key) {
		return "", fmt.Errorf("starts with a slash")
	}
	if key == "." {
		return "", fmt.Errorf("invalid path \".\"")
	}
	if strings.HasPrefix(key, "../") {
		return "", fmt.Errorf("starts with \"../\"")
	}
	return filepath.FromSlash(key), nil
}

func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (driver.Reader, error) {
	relpath, err := resolvePath(key)
	if err != nil {
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	path := filepath.Join(b.dir, relpath)
	if strings.HasSuffix(path, attrsExt) {
		return nil, fmt.Errorf("open file blob %s: extension %q cannot be directly read", key, attrsExt)
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fileError{relpath: relpath, msg: err.Error(), kind: driver.NotFound}
		}
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	xa, err := getAttrs(path)
	if err != nil {
		return nil, fmt.Errorf("open file attributes %s: %v", key, err)
	}
	if length == 0 {
		return reader{
			size: info.Size(),
			xa:   xa,
		}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("open file blob %s: %v", key, err)
		}
	}
	r := io.Reader(f)
	if length > 0 {
		r = io.LimitReader(r, length)
	}
	return reader{
		r:    r,
		c:    f,
		size: info.Size(),
		xa:   xa,
	}, nil
}

type reader struct {
	r    io.Reader
	c    io.Closer
	size int64
	xa   xattrs
}

func (r reader) Read(p []byte) (int, error) {
	if r.r == nil {
		return 0, io.EOF
	}
	return r.r.Read(p)
}

func (r reader) Close() error {
	if r.c == nil {
		return nil
	}
	return r.c.Close()
}

func (r reader) Attrs() *driver.ObjectAttrs {
	return &driver.ObjectAttrs{
		Size:        r.size,
		ContentType: r.xa.ContentType,
	}
}

func (b *bucket) NewTypedWriter(ctx context.Context, key string, contentType string, opt *driver.WriterOptions) (driver.Writer, error) {
	relpath, err := resolvePath(key)
	if err != nil {
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	path := filepath.Join(b.dir, relpath)
	if strings.HasSuffix(path, attrsExt) {
		return nil, fmt.Errorf("open file blob %s: extension %q is reserved and cannot be used", key, attrsExt)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("open file blob %s: %v", key, err)
	}
	attrs := xattrs{
		ContentType: contentType,
	}
	return &writer{
		w:     f,
		path:  path,
		attrs: attrs,
	}, nil
}

type writer struct {
	w     io.WriteCloser
	path  string
	attrs xattrs
}

func (w writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w writer) Close() error {
	if err := setAttrs(w.path, w.attrs); err != nil {
		return fmt.Errorf("write blob attributes: %v", err)
	}
	return w.w.Close()
}

func (b *bucket) Delete(ctx context.Context, key string) error {
	relpath, err := resolvePath(key)
	if err != nil {
		return fmt.Errorf("delete file blob %s: %v", key, err)
	}
	path := filepath.Join(b.dir, relpath)
	if strings.HasSuffix(path, attrsExt) {
		return fmt.Errorf("delete file blob %s: extension %q cannot be directly deleted", key, attrsExt)
	}
	err = os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fileError{relpath: relpath, msg: err.Error(), kind: driver.NotFound}
		}
		return fmt.Errorf("delete file blob %s: %v", key, err)
	}
	if err = os.Remove(path + attrsExt); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file blob %s: %v", key, err)
	}
	return nil
}

type fileError struct {
	relpath, msg string
	kind         driver.ErrorKind
}

func (e fileError) Error() string {
	return fmt.Sprintf("fileblob: object %s: %v", e.relpath, e.msg)
}

func (e fileError) BlobError() driver.ErrorKind {
	return e.kind
}
