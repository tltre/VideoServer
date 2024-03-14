package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatch send: %v", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
	forloop:
		for {
			select {
			case i := <-dc:
				log.Printf("Executor receive: %v", i)
			default:
				break forloop
			}
		}
		return nil
		// return errors.New("e")
	}

	runner := NewRunner(30, false, d, e)
	go runner.startAll()
	time.Sleep(3 * time.Second)
}
