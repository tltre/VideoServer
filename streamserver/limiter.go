package main

import "log"

type ConnLimiter struct {
	conCurrentConn int
	bucket         chan int
}

// NewConnLimiter 构造函数
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		conCurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

// GetConn 获取token
func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.conCurrentConn {
		log.Print("Reached the rate limitation")
		return false
	}

	log.Print("got a token")
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("New Connection Release: %d", c)
}
