package main

import (
	"log"
	"time"

	"github.com/Min-Feng/Channel_Application/golimit"
)

var counter int

func work(i int) {
	//log.Println("conter=",i)
	time.Sleep(time.Second)
}

func main(){
	maxGoroutine := 500
	limiter:= golimit.New(maxGoroutine)

	t1:=time.Now()
	for i:=0 ; i < 10000 ; i++{
		limiter.Add()  // 若goroutine數量超過限制,阻塞發生在Add
		taskID := i
		go func(){
			work(taskID)
			limiter.Done()
		}()
	}
	limiter.Wait()
	t2:=time.Since(t1)
	log.Println(t2)
}