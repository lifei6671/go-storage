package qiniu

import (
	"testing"
	"bytes"
	"github.com/lifei6671/go-storage"
)

var (
	key 	= ""
	secret 	= ""
)
func TestQiniuStorage_Write(t *testing.T) {

	ctx := storage.Context{}

	qiniu := NewQiniuStorage(key,secret,ctx)

	config := storage.Context{}

	config.Set("bucket","librarys")

	l,err := qiniu.Write("index.html","<html><body>你好</body></html>",config)

	if err != nil  {
		t.Fatal(err)
	}else{
		t.Log(l)
	}
}

func TestQiniuStorage_WriteStream(t *testing.T) {

	ctx := storage.Context{}

	qiniu := NewQiniuStorage(key,secret,ctx)

	config := storage.Context{}

	config.Set("bucket","librarys")

	b := []byte("微软股价受财报提振 周五大涨约5%\r\n")

	dstPath := "go/go-storage/TestQiniuStorage_WriteStream.txt"

	r := bytes.NewReader(b)

	err := qiniu.WriteStream(dstPath,r,int64(len(b)),config)

	if err != nil {
		t.Error(err)
	}
}
