package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

type Pool struct {
	mu        sync.Mutex
	resources chan io.Closer
	closed    bool
	factory   func() (io.Closer, error)
}

var ErrPoolClosed = errors.New("pool have been closed")

func New(fn func() (io.Closer, error), size int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value is negative")
	}
	return &Pool{
		resources: make(chan io.Closer, size),
		factory:   fn,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		log.Println("Acquire:", "share resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire:", "new resource")
		return p.factory()
	}
}

func (p *Pool) Release(r io.Closer) {
	// 釋放資源回到channel ,避免併發時
	// 其他gorountine關閉pool ,但Release動作沒有發現已經關閉
	// 向一個已經關閉的channel發送資料 ,會引起panic
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.resources <- r:
		log.Println("resource recycle to channel")
	default:
		log.Println("resources channel have been filled")
		r.Close()
	}
}

func (p *Pool) Close() {
	// 為了保證Release動作的安全性
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed == true {
		return
	}
	p.closed = true
	close(p.resources)

	// 注意: 關閉pool後,也要記得把channel中的資源釋放
	for r := range p.resources {
		r.Close()
	}
}
