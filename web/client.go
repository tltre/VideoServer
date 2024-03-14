package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

// Request API透传模块
func Request(apiBody *ApiBody, writer http.ResponseWriter, request *http.Request) {
	var resp *http.Response
	var err error

	// switch在匹配成功执行代码块后会自动退出，因此可以不用加break语句
	switch apiBody.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", apiBody.Url, nil)
		req.Header = request.Header
		resp, err = client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		normalResponse(writer, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", apiBody.Url, bytes.NewBuffer([]byte(apiBody.ReqBody)))
		req.Header = request.Header
		resp, err = client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		normalResponse(writer, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", apiBody.Url, nil)
		req.Header = request.Header
		resp, err = client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		normalResponse(writer, resp)
	default:
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "bad api request")
	}

}

func normalResponse(writer http.ResponseWriter, response *http.Response) {
	res, err := io.ReadAll(response.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		writer.WriteHeader(500)
		io.WriteString(writer, string(re))
		return
	}

	writer.WriteHeader(response.StatusCode)
	io.WriteString(writer, string(res))
}
