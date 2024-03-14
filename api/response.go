package main

import (
	"encoding/json"
	"io"
	"net/http"
	"src/api/defs"
)

/* -------- 处理返回响应 --------*/

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	res, _ := json.Marshal(errResp.Error)
	io.WriteString(w, string(res))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
