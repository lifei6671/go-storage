package blob

import (
	"sync"
	"errors"
	"sync/atomic"
	"fmt"
)

type atomicError struct{ v atomic.Value }

func (a *atomicError) Store(err error) {
	a.v.Store(struct{ error }{err})
}
func (a *atomicError) Load() error {
	err, _ := a.v.Load().(struct{ error })
	return err.error
}

var ErrClosedBuffer = errors.New("io: read/write on closed buffer")
var ErrBufferNotFull        = errors.New("bufio: buffer not full")

type BufferWriterReader struct {
	wrMu 	sync.Mutex // Serializes Write operations
	buf 	[]byte
	n 		int
	size 	int
	wrCh chan []byte

	once sync.Once // Protects closing done
	done chan struct{}
	rerr atomicError
	werr atomicError
}
func (p *BufferWriterReader) readCloseError() error {
	rerr := p.rerr.Load()
	if werr := p.werr.Load(); rerr == nil && werr != nil {
		return werr
	}
	return ErrClosedBuffer
}

func (p *BufferWriterReader) writeCloseError() error {
	werr := p.werr.Load()
	if rerr := p.rerr.Load(); werr == nil && rerr != nil {
		return rerr
	}
	return ErrClosedBuffer
}

//从缓冲区中读取一个数组，该数组长度可能会小于设定的缓冲长度
func (p *BufferWriterReader) ReadBytes() (buf []byte,err error)  {

	select {
	case <-p.done:
		return nil, p.readCloseError()
	default:
	}

	select {
	case b := <- p.wrCh:
		err = nil
		buf = b
		return
	case <-p.done:
		if p.n > 0 {
			buf = p.buf
			p.buf = make([]byte,p.size)
			p.n = 0
			return
		}
		return nil, p.readCloseError()
	default:
		return nil,ErrBufferNotFull
	}
	return nil,ErrBufferNotFull
}

func (p *BufferWriterReader) Write(b []byte) (n int, err error) {
	select {
	case <-p.done:
		return 0, p.writeCloseError()
	default:
		p.wrMu.Lock()
		defer p.wrMu.Unlock()
	}


	for once := true; once || len(b) > 0; once = false {
		l :=  len(p.buf) - p.n

		fmt.Println("bbbb",l)
		var j int
		if len(b) > l {
			j = copy(p.buf[p.n:], b[0:l])
			b = b[l+1:]
		}else{
			j = copy(p.buf[p.n:],b)

			fmt.Println("aaa",j)

			b = b[:0]
		}

		p.n += j
		n += j


		if p.n == p.size {

			select {
			case p.wrCh <- p.buf:
				p.buf = make([]byte,p.n)
				p.n = 0
				fmt.Println("ok")
				break
			case <-p.done:
				return n, p.writeCloseError()
			}
		}
		if len(b) > 0 {
			once = true
		}
	}
	return n, nil
}

func (b *BufferWriterReader) Close() error {

	b.once.Do(func() { close(b.done) })

	return nil

}

func NewBufferWriterReaderSize(rSize int,wSize int) *BufferWriterReader {
	return &BufferWriterReader{
		wrMu: sync.Mutex{},
		buf: make([]byte,rSize),
		wrCh: make(chan []byte,wSize),
		n:0,
		size:rSize,
	}
}