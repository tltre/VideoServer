package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"src/api/dbops"
	"src/api/defs"
	"src/api/session"
	"src/api/utils"
)

/* ------------ User Handler ------------ */

// CreateUser 新建用户
func CreateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	res, _ := io.ReadAll(request.Body)
	ubody := &defs.UserCredential{}

	// 将json转成结构体
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := defs.SignedUp{
		Success:   true,
		SessionId: id,
	}
	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 201)
	}
}

// Login 用户登录
func Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// body中存有用户名与密码
	res, _ := io.ReadAll(request.Body)
	log.Print(res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}

	// 在查询行参数中获得username
	userName := params.ByName("user_name")
	log.Printf("Url username: %s", userName)
	log.Printf("RequestBody username: %s", ubody.Username)
	if userName != ubody.Username {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	pwd, err := dbops.GetUserCredential(userName)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}
	log.Printf("Login pwd: %s", pwd)
	log.Printf("RequestBody pwd: %s", ubody.Pwd)

	// 每次登录返回一个新的sessionId
	id := session.GenerateNewSessionId(userName)
	su := defs.SignedIn{
		Success:   true,
		SessionId: id,
	}
	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 202)
	}
}

// GetUserInfo 获取用户信息
func GetUserInfo(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// 校验是否已登录
	if !validateUser(writer, request) {
		log.Print("Unauthorized User\n")
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	userName := params.ByName("user_name")
	id, err := dbops.GetUserId(userName)
	if err != nil {
		log.Printf("Error in GetUserInfo: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	user := &defs.UserInfo{
		Username: userName,
		UserId:   id,
	}

	if resp, err := json.Marshal(user); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 200)
	}
}

// DeleteUser 注销用户
func DeleteUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// 查验是否通过健全验证（session是否已过期，若过期则提示登录
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	// 参数行中放入用户名与密码
	userName := params.ByName("user_name")
	pwd := params.ByName("pwd")

	// 验证是否有该用户
	DBpwd, err := dbops.GetUserCredential(userName)
	if err != nil {
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	if DBpwd == "" || DBpwd != pwd {
		// 无该用户
		sendErrorResponse(writer, defs.ErrorNoSuchUser)
		return
	}

	// 删除user表中用户
	err = dbops.DeleteUser(userName, pwd)
	if err != nil {
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	// 删除该用户的session
	sid := request.Header.Get(HeaderFieldSession)
	err = session.DeleteExpiredSession(sid)
	if err != nil {
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	// 注意响应码问题：使用204时是 No Content，因此会不显示出回答
	sendNormalResponse(writer, "Successfully DELETE", 200)
}

/* ------------ Video Handler ------------ */

func PostNewVideo(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// 查验是否通过健全验证（session是否已过期，若过期则提示登录
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	res, _ := io.ReadAll(request.Body)
	nvBody := &defs.NewVideoInfo{}
	if err := json.Unmarshal(res, nvBody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(nvBody.AuthorId, nvBody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 200)
	}
}

func ListAllVideos(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// 查验是否通过健全验证（session是否已过期，若过期则提示登录
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	username := params.ByName("username")
	vs, err := dbops.ListVideo(username, 0, utils.GetCurrentTimestampSec())

	if err != nil {
		log.Printf("Error in ListVideos: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	videos := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(videos); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 200)
	}
}

func DeleteVideo(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// 查验是否通过健全验证（session是否已过期，若过期则提示登录
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	vid := params.ByName("vid-id")
	err := dbops.DeleteVideo(vid)
	if err != nil {
		log.Printf("Error in DeleteVideo: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	sendNormalResponse(writer, "", 204)
}

/* ------------ Comment Handler ------------ */

func PostNewComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	res, _ := io.ReadAll(request.Body)
	cBody := &defs.NewComment{}
	if err := json.Unmarshal(res, cBody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(writer, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := params.ByName("vid-id")
	err := dbops.AddNewComment(cBody.AuthorId, vid, cBody.Content)
	if err != nil {
		log.Printf("Error in PostComments: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	} else {
		sendNormalResponse(writer, "ok", 200)
	}
}

func ShowComments(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !validateUser(writer, request) {
		sendErrorResponse(writer, defs.ErrorNotAuthUser)
		return
	}

	vid := params.ByName("vid-id")
	res, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(writer, defs.ErrorDBError)
		return
	}

	comments := &defs.CommentList{Comments: res}
	if resp, err := json.Marshal(comments); err != nil {
		sendErrorResponse(writer, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(writer, string(resp), 200)
	}
}
