package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func testPageHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w,nil)
}

func indexPageHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	t, _ := template.ParseFiles("./videos/video.html")
	t.Execute(w,nil)
}


func streamHandler(w http.ResponseWriter,r *http.Request, p httprouter.Params)  {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid
	video , err :=os.Open(vl)
	if err != nil {
		log.Printf("internal error")
		sendErrorResponse(w,http.StatusInternalServerError,"internal error")
		return
	}
	w.Header().Set("Content-type","video/mp4")
	http.ServeContent(w,r,"",time.Now(),video)
	defer video.Close()
}

func uploadHandler(w http.ResponseWriter,r *http.Request, p httprouter.Params)  {
	r.Body = http.MaxBytesReader(w,r.Body,MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil{
		sendErrorResponse(w,http.StatusBadRequest,"File is too big")
		return
	}
	file,_,err := r.FormFile("file")
	if err != nil{
		sendErrorResponse(w,http.StatusInternalServerError,"interner error")
		return
	}
	data ,err := ioutil.ReadAll(file)
	if err != nil{
		log.Printf("Read file error: %v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"Read file error")
		return
	}
	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn,data,0666)
	if err != nil{
		log.Printf("write file error: %v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"write file error")
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w,"uploaded successfully")
}