package main

import (
	"sync"
	"log"
	"time"

	"github.com/Min-Feng/Channel_Application/work"
)

func testFunc(i int) {
	//log.Println("conter=",i)
	time.Sleep(time.Second)
}

type taskInt int

func (t *taskInt) Task(){
	testFunc(int(*t))
}

func main() {
	maxGoroutine := 500
	workPool:=work.New(maxGoroutine)
	var wg sync.WaitGroup

	t1:=time.Now()
	for i:= 0 ; i < 10000 ; i++{
		wg.Add(1)
		taskID := taskInt(i)
		go func(){
			workPool.Run(&taskID)   // 若goroutine數量超過限制,阻塞發生在Run
			wg.Done()
		}()
	}
	wg.Wait()
	workPool.Close()  // Close()方法,直到所有工作都完成才會返回
	t2:=time.Since(t1)
	log.Println(t2)
}