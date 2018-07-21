package local

import (
	"os"
	"time"
	"mime"
	"path"
)

type LocalFileInfo struct {
	f os.FileInfo
	hash string
	path string
	mimeType string
}

func NewLocalFileInfo(f os.FileInfo) *LocalFileInfo {
	return &LocalFileInfo{f: f  }
}

func (f *LocalFileInfo ) Name() string {
	return f.f.Name()
}
func (f *LocalFileInfo )Size() int64 {
	return f.f.Size()
}
func (f *LocalFileInfo )ModTime() time.Time {
	return f.f.ModTime()
}
func (f *LocalFileInfo )MimeType() string {
	if f.mimeType == "" {

		mimetype := mime.TypeByExtension(path.Ext(f.path))

		f.mimeType = mimetype
	}
	return f.f.Name()
}
func (f *LocalFileInfo )Hash() string {
	if f.hash == "" {

	}
	return f.hash
}
func (f *LocalFileInfo)IsDir() bool  {
	return f.f.IsDir()
}
func (f *LocalFileInfo)Sys() interface{}   {
	return f.f.Sys()
}