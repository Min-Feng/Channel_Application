package main

import (
	"log"
	"os"
	"time"

	"github.com/Min-Feng/concurrencyPractice/runner"
)

const timeout = 3 * time.Second

func main() {
	log.Println("starting work")

	r := runner.New()

	f := createTask()
	r.Add(f, f, f)

	if err := r.Start(timeout); err != nil {
		switch err {
		case runner.ErrTimeOut:
			log.Println(err)
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println(err)
			os.Exit(2)
		}
	}
	log.Println("tasks finish")
}

func createTask() func() {
	var id int
	return func() {
		id++
		log.Printf("processor Task #%d\n", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
