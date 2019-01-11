package pool

import (
	"io"
	"sync"
	"errors"
)

// Singleton Pattern 單例模式
// 保證只有一個實例
var p Pool

type Pool struct {
	m        sync.Mutex
	resource chan io.Closer
	closed   bool
	factory  func() (io.Closer, error)
}

var ErrPoolClosed=errors.New("pool have been closed.")


