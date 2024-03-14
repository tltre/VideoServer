package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !m.l.GetConn() {
		sendErrorResponse(writer, http.StatusTooManyRequests, "To many requests")
		return
	}
	m.r.ServeHTTP(writer, request)
	defer m.l.ReleaseConn()
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()

	// 获取视频
	router.GET("/videos/:vid-id", streamHandler)

	// 上传视频
	router.POST("/upload/:vid-id", uploadHandler)

	// 测试
	router.GET("/testPage", testPageHandler)

	return router
}

func main() {
	r := RegisterRouter()
	m := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", m)
}
