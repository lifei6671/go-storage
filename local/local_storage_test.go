package local

import (
	"testing"
	"os"
	"path/filepath"
	"bytes"
	"log"
)

var dstPath = filepath.Join(os.TempDir(),"go-storage","local_storage.test.text")

func TestLocalStorage_Stat(t *testing.T) {

	localStorage := NewLocalStorage()

	f,err := localStorage.Stat(dstPath,nil)

	if err != nil {
		t.Error(err)
	}else{
		t.Log(f.Name(),f.ModTime(),f.Size())
	}
}

func TestLocalStorage_Write(t *testing.T) {
	localStorage := NewLocalStorage()

	i,err := localStorage.Write(dstPath,"微软股价受财报提振 周五大涨约5%",nil)

	if err != nil {
		t.Error(err)
	}else{
		t.Log(dstPath,"->",i)
	}
}

func TestLocalStorage_WriteBytes(t *testing.T) {
	localStorage := NewLocalStorage()

	i,err := localStorage.WriteBytes(dstPath,[]byte("微软股价受财报提振 周五大涨约5%"),nil)

	if err != nil {
		t.Error(err)
	}else{
		t.Log(dstPath,"->",i)
	}
}

func TestLocalStorage_WriteStream(t *testing.T) {
	localStorage := NewLocalStorage()

	b := []byte("微软股价受财报提振 周五大涨约5%\r\n")

	r := bytes.NewReader(b)

	err := localStorage.WriteStream(dstPath,r,int64(len(b)),nil)

	if err != nil {
		t.Error(err)
	}
}

func TestLocalStorage_Append(t *testing.T) {
	localStorage := NewLocalStorage()
	l,err := localStorage.Append(dstPath,"如何看待锤子科技官网取消TNT全款预订，并且罗永浩微博删除所有TNT相关微博？",nil)

	if err != nil {
		t.Error(err)
	}else{
		t.Log(dstPath,"->",l)
	}

}

func TestLocalStorage_AppendBytes(t *testing.T) {
	localStorage := NewLocalStorage()
	l,err := localStorage.AppendBytes(dstPath,[]byte("如何看待锤子科技官网取消TNT全款预订，并且罗永浩微博删除所有TNT相关微博？"),nil)

	if err != nil {
		t.Error(err)
	}else{
		t.Log(dstPath,"->",l)
	}
}

func TestLocalStorage_AppendStream(t *testing.T) {
	localStorage := NewLocalStorage()

	b := []byte("微软股价受财报提振 周五大涨约5%\r\n")

	r := bytes.NewReader(b)

	err := localStorage.AppendStream(dstPath,r,int64(len(b)),nil)

	if err != nil {
		log.Fatal(err)
	}
}