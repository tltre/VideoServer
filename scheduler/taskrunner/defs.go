package taskrunner

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE  = "e"
	CLOSE             = "c"

	VIDEO_PATH = ".\\videos\\"
)

// control Channel连接生产者与消费者，指明能否进行处理
type controlChan chan string

// data Channel表示下发的数据
// 使用interface表示各种类型
type dataChan chan interface{}

// fn 是 dispatcher（生产者） 和 executor（消费者）
type fn func(dc dataChan) error
