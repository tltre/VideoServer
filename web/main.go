package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	// 未登录界面，第一次是Get方法，第二次提交表单登录是Post方法
	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	// 用户界面，同样需要有Get和Post方法
	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	// API透传方式调用后端API
	router.POST("/api", apiHandler)

	// 对于一些API调用（如无法规范化成APIBody形式）
	// 为了避免跨域访问问题，使用代理模式调用API
	router.POST("/upload/:vid-id", proxyHandler)

	// 设置静态资源路径
	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8080", r)
}
