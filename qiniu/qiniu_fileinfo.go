package qiniu

import (
	"time"
	"github.com/qiniu/api.v7/storage"
)

//七牛文件信息
type QiniuFileInfo struct {
	info storage.FileInfo
	name string
}

func NewQiniuFileInfo(info storage.FileInfo) *QiniuFileInfo {
	return &QiniuFileInfo{ info:info}
}

func (f *QiniuFileInfo) Name() string {
	return f.name
}
func (f *QiniuFileInfo) Size() int64{
	return f.info.Fsize
}
func (f *QiniuFileInfo)ModTime() time.Time{
	 return  time.Unix(f.info.PutTime,0)
}
func (f *QiniuFileInfo) MimeType() string{
	return f.info.MimeType
}
func (f *QiniuFileInfo)Hash() string{
	return f.info.Hash
}
func (f *QiniuFileInfo)IsDir() bool  {
	return false
}
func (f *QiniuFileInfo)Sys() interface{}   {
	return nil
}