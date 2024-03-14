package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid + ".mp4"

	video, err := os.Open(vl)
	if err != nil {
		log.Printf("error when try to open file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// 强制设置content-type为mp4
	w.Header().Set("Content-Type", "video/mp4")
	// 设置播放方式
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 检查文件大小是否超出限制
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}

	// 读取 <form name = "file"> 内的文件
	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	// 将文件的二进制流写到server端的新文件内
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error of read file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	fn := p.ByName("vid-id")
	err = os.WriteFile(VIDEO_DIR+fn+".mp4", data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded Successfully")

	defer file.Close()

}
