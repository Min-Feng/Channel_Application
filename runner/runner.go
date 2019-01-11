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
	r.timeout = time.After(d)
	go func() {
		//因為會先執行main gorountine,timeout變數還是舊的,因此永遠不會接收到終止訊號
		//儘量將time.After的發送和接收處於同一個goroutine
		//r.timeout = time.After(d)  
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout: // 無論任務是否進行中,只要時間超過,就結束
		return ErrTimeOut
	}

}

func (r *Runner) run() error {
	for _, task := range r.tasks {
		select {
		case <-r.interrupt:
			// 停止接收後續的OS信號,只有在每個任務的間距,才會執行中斷動作
			signal.Stop(r.interrupt)
			return ErrInterrupt
		default:
			task()
		}
	}
	return nil
}
