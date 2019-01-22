
// 針對特定函數,限制同時最多幾個人使用
package main

import (
	"log"
	"sync"
)

func work(i int) {
	log.Println("work num", i)
}

func funcPool(fn func(int), size int) chan func(int) {
	ch := make(chan func(int), size)
	for i := 0; i < size; i++ {
		ch <- fn
	}
	return ch
}

var counterCalling int

func run(ch chan func(int)) {
	var mu sync.Mutex
	select {
	case fn := <-ch:
		mu.Lock()
		fn(counterCalling)
		counterCalling++
		mu.Unlock()
		// ch <- fn
	default:
		log.Println("超過最大使用人數,稍後再呼叫")
	}
}

func main() {
	ch := funcPool(work, 10)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			run(ch)
			wg.Done()
		}()
	}
	wg.Wait()
}
