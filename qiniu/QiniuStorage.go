package qiniu

import (
	"io"
	"github.com/lifei6671/mindoc/storage"
)

type QiniuStorage struct {

}




func (s *QiniuStorage) Stat(path string,context storage.Context) (storage.FileInfo,error){

}
func (s *QiniuStorage) Write(path string, contents string, context storage.Context) (int64,error){

}
func (s *QiniuStorage) WriteStream(path string, reader io.Reader, context storage.Context) (int64,error){

}
func (s *QiniuStorage) WriteBytes(path string, contents []byte, context storage.Context) (int64,error){

}
func (s *QiniuStorage) Append(path string, contents string, context storage.Context) (int64,error){

}
func (s *QiniuStorage) AppendStream(path string, reader io.Reader, context storage.Context) (int64,error){

}
func (s *QiniuStorage) AppendBytes(path string, contents []byte, context storage.Context) (int64,error){

}
func (s *QiniuStorage) ReName(path string,newPath string,context storage.Context) error{

}
func (s *QiniuStorage) Copy (src string,dst string, context storage.Context) error{

}
func (s *QiniuStorage) Delete(dst string,context storage.Context) error{

}
func (s *QiniuStorage) DeleteDir(dst string,context storage.Context) error{

}
func (s *QiniuStorage) CreateDir(dst string, context storage.Context) error{

}

