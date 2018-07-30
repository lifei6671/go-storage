package blob

import (
	"testing"
	"fmt"
)

func TestBufferWriterReader_Write(t *testing.T) {
	rw := NewBufferWriterReaderSize(1024,1)

	buf := make([]byte,1024)

	for i := 0; i < 1024;i++ {
		buf[i] = byte('b')
	}

		l,err := rw.Write(buf)

		if err != nil {
			t.Fatal(err)
		}
		t.Log(l)


	b,err := rw.ReadBytes()

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
