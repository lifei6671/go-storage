package qiniublob

import (
	"testing"
	"context"
	"github.com/qiniu/api.v7/storage"
	"os"
	"io/ioutil"
	"io"
	"fmt"
	"github.com/lifei6671/go-storage/blob"
)

var (
	key 	= "LY5mwlQ3vVqfof3FWCionaDQuqz3sM4N6aqaQ8wQ"
	secret 	= "woXMxDAng6beeqCAhdDRxAng3lPVlDj88QEVFJRa"
)

func TestBucket_NewRangeReader(t *testing.T) {
	cfg := storage.Config{
		Zone: &storage.ZoneHuadong,
		UseHTTPS:false,
	}

	bucket,err := OpenBucket(key,secret,cfg,"librarys","http://libs.iminho.me")

	if err != nil {
		t.Fatal(err)
	}
	key := "d7c3bd4cly1fsawnkpo91j23vc2kw4qt.jpg"

	reader,err := bucket.NewRangeReader(context.Background(),key,0,-1)

	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	f,err := ioutil.TempFile(os.TempDir(),key)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	_,err = io.Copy(f,reader)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f.Name(), reader.ContentType(),reader.Size())
	t.Log("文件为：",f.Name())

}

func TestBucket_NewTypedWriter(t *testing.T) {
	p := "index.htm"

	cfg := storage.Config{
		Zone: &storage.ZoneHuadong,
		UseHTTPS:false,
	}

	bucket,err := OpenBucket(key,secret,cfg,"librarys","http://libs.iminho.me")

	if err != nil {
		t.Fatal(err)
	}

	writer,err := bucket.NewWriter(context.Background(),p,&blob.WriterOptions{ BufferSize:0,ContentType:"text/html"})

	if err != nil {
		t.Fatal(err)
	}
	l,err := writer.Write([]byte("<html><body>你好</body></html>"))
	if err != nil {
		t.Fatal(err,l)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err,l)
	}
	if l <= 0 {
		t.Fatal("上传文件失败")
	}
}

func TestBucket_NewTypedWriter2(t *testing.T) {
	p := "006pdpXYly1flasoa1xgmj31ww2gh4qs.jpg"

	cfg := storage.Config{
		Zone: &storage.ZoneHuadong,
		UseHTTPS:false,
	}

	bucket,err := OpenBucket(key,secret,cfg,"librarys","http://libs.iminho.me")

	if err != nil {
		t.Fatal(err)
	}

	writer,err := bucket.NewWriter(context.Background(),p,&blob.WriterOptions{ BufferSize:1024*20, ContentType:"text/html"})

	if err != nil {
		t.Fatal(err)
	}
	b,err := ioutil.ReadFile("/Users/minho/OneDrive/图片/006pdpXYly1flasoa1xgmj31ww2gh4qs.jpg")

	if err != nil {
		t.Fatal(err)
	}

	l,err := writer.Write(b)
	if err != nil {
		t.Fatal(err,l)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err,l)
	}
	if l <= 0 {
		t.Fatal("上传文件失败")
	}
}