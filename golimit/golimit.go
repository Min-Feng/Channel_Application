package golimit

import (
	"sync"
)

type Limit struct {
	queue chan struct{}
	wg    *sync.WaitGroup
}

const (
	defaultQueueSize = 500
)

// New 生成限制器實例,內部維護一個channel,以便於控制goroutine的數量在指定的範圍內
func New(size int) *Limit {
	if size <= 0 {
		size = defaultQueueSize
	}

	return &Limit{
		queue: make(chan struct{}, size),
		wg:    new(sync.WaitGroup),
	}
}

// Add 確認一個goroutine已經開始執行,若列隊已滿,則阻塞,確保goroutine維持在一定數量
func (l *Limit) Add() {
	l.queue <- struct{}{}
	l.wg.Add(1)
}

// Done 若goroutine已完成,則退出列隊,以便於其他的goroutine進入列隊
func (l *Limit) Done() {
	<-l.queue
	l.wg.Done()
}

// Wait 等待所有goroutine完成
func (l *Limit) Wait() {
	l.wg.Wait()
}
