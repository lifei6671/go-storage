package storage

import (
	"io"
)

type Adapter interface {
	Stat(path string,context Context) (FileInfo,error)
	Write(path string, contents string, context Context) (int64,error)
	WriteStream(path string, reader io.Reader, context Context) (int64,error)
	WriteBytes(path string, contents []byte, context Context) (int64,error)
	Append(path string, contents string, context Context) (int64,error)
	AppendStream(path string, reader io.Reader, context Context) (int64,error)
	AppendBytes(path string, contents []byte, context Context) (int64,error)
	ReName(path string,newPath string,context Context) error
	Copy (src string,dst string, context Context) error
	Delete(dst string,context Context) error
	DeleteDir(dst string,context Context) error
	CreateDir(dst string, context Context) error
}
