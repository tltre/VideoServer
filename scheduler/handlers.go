package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"src/scheduler/dbops"
)

func VideoDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")

	if len(vid) == 0 {
		sendResponse(w, 400, "video id should not be empty")
		return
	}

	err := dbops.AddVideoDelRecord(vid)
	if err != nil {
		sendResponse(w, 500, "Internal error")
		return
	}

	sendResponse(w, 200, "success")
	return
}
