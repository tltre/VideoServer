package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"src/scheduler/taskrunner"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", VideoDelRecHandler)

	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	// 该函数已经阻塞
	http.ListenAndServe(":9001", r)
}
