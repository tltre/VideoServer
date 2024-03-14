package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// cookie中需要放入username和session字段
	cname, err1 := request.Cookie("username")
	sid, err2 := request.Cookie("session")

	// 出错进入登录界面
	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "xy"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing templates home.html error : %s", e)
			return
		}
		t.Execute(writer, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		// 若sid符合，则重定向
		// 判断是否符合的工作放在前端，所以这里只需要存在即可
		http.Redirect(writer, request, "/userhome", http.StatusFound)
	}
}

func userHomeHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cname, err1 := request.Cookie("username")
	_, err2 := request.Cookie("session")

	// 出错进入登录界面
	if err1 != nil || err2 != nil {
		http.Redirect(writer, request, "/", http.StatusFound)
		return
	}

	// 表单提交的用户名，合法性在后端验证
	fname := request.FormValue("username")

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Error when parsing html file: %s", e)
		return
	}

	t.Execute(writer, p)
}

// 使用API透传方式处理前台发过来的调用API请求
func apiHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if request.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognize)
		io.WriteString(writer, string(re))
		return
	}

	res, _ := io.ReadAll(request.Body)
	apiBody := &ApiBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(writer, string(re))
		return
	}

	Request(apiBody, writer, request)

	defer request.Body.Close()
}

// 使用代理方式处理前台发过来的调用API请求
func proxyHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	u, _ := url.Parse(StreamServerUrl)
	// proxy将原本的域名替换为目标域名
	// 并不会修改头部内容
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(writer, request)
}
