package qiniu

import (
	"time"
	"github.com/qiniu/api.v7/storage"
)

type QiniuFileInfo struct {
	info storage.FileInfo
}

func (f *QiniuFileInfo) Name() string {
	return f.name
}
func (f *QiniuFileInfo) Size() int64{
	return 0
}
func (f *QiniuFileInfo)ModTime() time.Time{
	return time.Now()
}
func (f *QiniuFileInfo) MimeType() string{
	return ""
}
func (f *QiniuFileInfo)Hash() string{
	return ""
}
