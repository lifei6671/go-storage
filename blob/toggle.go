package blob

import (
	"sync/atomic"
	"sync"
)

// 第一次执行f1之后一直执行f2.
type Toggle struct {
	m    sync.Mutex
	done uint32
}

func (o *Toggle) Do(f1 func(),f2 func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		f2()
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f1()
	}
}

