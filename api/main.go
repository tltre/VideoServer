package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"src/api/session"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func (m middleWareHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// check session
	validateUserSession(request)
	m.r.ServeHTTP(writer, request)
}

func NewMiddleWare(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func RegistersHandlers() *httprouter.Router {
	router := httprouter.New()

	// 用户注册（无参数
	router.POST("/user", CreateUser)

	// 用户登录（传入参数
	router.POST("/user/:user_name", Login)

	// 获取用户信息
	router.GET("/user/:user_name", GetUserInfo)

	// 用户注销
	router.DELETE("/user/:user_name/:pwd", DeleteUser)

	router.POST("/user/:username/videos", PostNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostNewComment)

	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegistersHandlers()
	mh := NewMiddleWare(r)
	http.ListenAndServe(":8000", mh)
}
