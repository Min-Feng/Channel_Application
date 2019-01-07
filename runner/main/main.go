package main

import(
	"log"
	"os"
	"time"

	"github.com/Min-Feng/concurrencyPractice/runner"
)

const timeout = 3* time.Second

func main () {
	log.Println("starting work")

	r:= runner.New(timeout)

	r.Add(createTask(),createTask(),createTask())

	if err:=r.Start(); err != nil{
		switch err{
		case runner.ErrTimeOut:
			log.Println(err)
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println(err)
			os.Exit(2)
		}
	}

	f:=createTask()
	f()
	f()
}

func createTask() func() {
	var id int
	id++
	return func(){
		log.Printf("processor Task #%d\n",id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}