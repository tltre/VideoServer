package taskrunner

import (
	"errors"
	"log"
	"src/scheduler/dbops"
	"sync"
)

/* ------ 生产者-消费者任务 ------ */
/* -------- 延时删除任务 -------- */

// VideoClearDispatch 将待删除的video_id从数据库中读到管道中等待消费者处理
func VideoClearDispatch(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatch error: %v", err)
	}
	if len(res) == 0 {
		return errors.New("all tasks finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
forloop:
	for {
		select {
		case vid := <-dc:
			// 使用go func会有重复读写问题：也就是还未删除完成就通知生产者可以生产，生产者会又去取到之前的未删除的记录
			// 若不以id进行传参，会使用实时vid导致错误
			go func(id interface{}) {
				if err := dbops.DeleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}
