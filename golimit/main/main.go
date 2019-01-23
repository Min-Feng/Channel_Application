package main

import (
	"log"
	"time"

	"github.com/Min-Feng/Channel_Application/golimit"
)

var counter int

func work(i int) {
	log.Println("conter=",i)
	time.Sleep(time.Second)
}

func main(){
	limiter:= golimit.New(500)

	t1:=time.Now()
	for i:=0 ; i < 1000 ; i++{
		limiter.Add()
		go func(i int){
			work(i)
			limiter.Done()
		}(i)
	}
	limiter.Wait()
	t2:=time.Since(t1)
	log.Println(t2)
}