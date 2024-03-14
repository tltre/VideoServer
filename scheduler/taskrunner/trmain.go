package taskrunner

import "time"

// Worker 定时器与生产者-消费者模型的组装结合
type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) StartWorker() {
	for {
		select {
		// 若用range，无数据到来时会阻塞
		case <-w.ticker.C:
			go w.runner.startAll()
		}
	}
}

func Start() {
	// start video file cleaning
	r := NewRunner(3, true, VideoClearDispatch, VideoClearExecutor)
	w := NewWorker(3, r)
	w.StartWorker()
}
