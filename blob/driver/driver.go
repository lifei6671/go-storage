package driver

import "io"
import "context"

//源文件 https://github.com/google/go-cloud/blob/master/blob/driver/driver.go

// 错误标识码.
type ErrorKind int

const (
	GenericError ErrorKind = iota
	NotFound
)

// 封装的错误接口
type Error interface {
	error
	BlobError() ErrorKind
}

type Reader interface {
	io.ReadCloser
	Attrs() *ObjectAttrs
}

// 写接口.
type Writer interface {
	io.WriteCloser
}

// 控制写的选项.
type WriterOptions struct {
	//分片写入大小限制，如果小于设置的大小，将单次请求写入，否则在服务器端支持的情况下，分片多次请求写到服务器端
	BufferSize int
}

// 对象包含的元数据.
type ObjectAttrs struct {
	// 对象大小.
	Size int64
	// ContentType 对象的MIME类型，不能为空.
	ContentType string
}

type Bucket interface {
	//读取给定对象指定的偏移量的一部分，如果偏移量为0则只返回元数据，如果为负数，则返回对象的所有数据
	NewRangeReader(ctx context.Context, key string, offset, length int64) (Reader, error)
	//将指定的对象写入到服务器端，如果对象不存在则创建新对象，如果已存在则覆盖，写入完成后需要手动调用Close方法关闭写入对象。对象类型不能为空
	NewTypedWriter(ctx context.Context, key string, contentType string, opt *WriterOptions) (Writer, error)
	// 删除与指定key管理的对象，如果对象不存在返回一个错误
	Delete(ctx context.Context, key string) error
}
