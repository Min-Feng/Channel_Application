package work

import (
	"sync"
)

type worker interface {
	Task()
}

type Pool struct {
	works chan worker
	wg    sync.WaitGroup
}

const (
	defaultMaxGoroutine = 500
)

// New 創建新的pool池的同時,也建立了所要求的goroutine數量
func New(maxGoroutine int) *Pool {
	if maxGoroutine <= 0 {
		maxGoroutine = defaultMaxGoroutine
	}

	p := new(Pool)
	p.works = make(chan worker)
	p.wg.Add(maxGoroutine)

	for i := 0; i < maxGoroutine; i++ {
		go func() {
			for w := range p.works {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return p
}

func (p *Pool) Run(w worker) {
	p.works <- w
}

func (p *Pool) Close() {
	close(p.works)
	p.wg.Wait()
}
