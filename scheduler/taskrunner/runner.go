package taskrunner

// Runner 生产者-消费者模型
type Runner struct {
	// 传输ready信息
	Controller controlChan

	// 传输close信息
	Error controlChan

	// 数据通道
	Data dataChan

	Dispatch fn
	Executor fn

	dataSize  int
	longLived bool
}

func NewRunner(size int, longLived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		Dispatch:   d,
		Executor:   e,
		dataSize:   size,
		longLived:  longLived,
	}
}

func (r *Runner) startDispatch() {
	// 根据longLived选择是否回收资源
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Error)
			close(r.Data)
		}
	}()

	for {
		// select 进行异步监听，Controller有内容则进入case c；Error有内容进入case e
		select {
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatch(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == CLOSE {
				return
			}
		default:
		}
	}
}

func (r *Runner) startAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
