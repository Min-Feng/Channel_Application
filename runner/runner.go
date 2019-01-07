package runner

// 使用channel來監視程序的執行時間,若執行時間太長,則終止程序
import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// Runner 在給定的一組時間內,進行一系列任務task
// 且在併發執行時,發送終止訊號來結束任務
type Runner struct {
	// os傳來的發送信號
	interrupt chan os.Signal

	// 報告任務是否完成
	complete chan error

	// 報告任務是否超時
	timeout <-chan time.Time

	// 依序執行各式函數
	tasks []func()
}

// ErrTimeOut 任務執行超時時返回
var ErrTimeOut = errors.New("received timeout")

// ErrInterrupt 接收到OS的事件時返回
var ErrInterrupt = errors.New("received interrupt")

// New 返回一個Runner實例
func New() *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
	}
}

// Add 將想執行的任務加入到Runner實例
func (r *Runner) Add(task ...func()) {
	r.tasks = append(r.tasks, task...)
}

// Start 執行Runner實例的任務,設定 d 時間內完成任務
func (r *Runner) Start(d time.Duration) error {
	// 接收os的中斷訊號,傳送到該chan
	signal.Notify(r.interrupt, os.Interrupt)

	go func(){
		r.complete <- r.run(d)
	}()

	err := <-r.complete
	return err
}

func (r *Runner) run(d time.Duration) error {
	r.timeout=time.After(d)
	
	for _,task := range r.tasks{
		select{
		case <- r.interrupt:
			// 停止接收後續的OS信號
			signal.Stop(r.interrupt)
			return ErrInterrupt
		case <-r.timeout:
			return ErrTimeOut
		default:
			task()
		}
	}
	return nil
}
