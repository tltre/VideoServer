package main

import (
	"net/http"
	"src/api/session"
)

/* -------- 健全验证等 --------*/

var HeaderFieldSession = "X-Session-Id"
var HeaderFieldUname = "X-User-Name"

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HeaderFieldSession)
	if len(sid) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HeaderFieldUname, uname)
	return true
}

func validateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HeaderFieldUname)
	if len(uname) == 0 {
		// sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}
