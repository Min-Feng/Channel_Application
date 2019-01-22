package main

import (
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/Min-Feng/Channel_Application/pool"
)

const (
	maxGoroutine     = 25
	poolResourceSize = 4
)

var idCounter int32

// dbConnection 模擬想重複使用的資源
type dbConnection struct {
	ID int32
}

func (dbConn *dbConnection) Close() error {
	log.Println("Close: Connection", dbConn.ID)
	return nil
}

func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create: New Connection", id)
	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutine)

	pool, err := pool.New(createConnection, poolResourceSize)
	if err != nil {
		log.Println("pool create failed", err)
	}

	for query := 0; query < maxGoroutine; query++ {
		// 瞬間的資源要求數量,若遠大於pool size的數量
		// pool的資源重複用目的,無法妥善使用
		// 因此選擇合適的pool size是很重要的
		go func(q int) {
			performQueries(q, pool)
			wg.Done()
		}(query)
	}

	wg.Wait()
	log.Println("shutdown program")
	pool.Close()
}

func performQueries(query int, p *pool.Pool) {
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	// 每次獲得資源 要記得釋放或關閉資源 !!!
	defer p.Release(conn)

	log.Printf("Query ID[%d] Conn ID[%d] \n",query,conn.(*dbConnection).ID)

}
