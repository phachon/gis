package main

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"os"
)

// go 实现的图片上传,存储,缩放,下载服务

func init()  {
	initConfig();
}

func main()  {

	go downloadServer()
	uploadServer()
}

// 上传 server
func uploadServer()  {

	uploadServer := conf.GetString("listen.upload")
	router := httprouter.New()
	httpHandle := NewHttpHandle()
	router.GET("/", httpHandle.Index)
	router.POST("/image/upload", httpHandle.ImageUpload)
	//跨域
	router.OPTIONS("/image/upload", httpHandle.CrossDomain)

	log.Println("upload server start listen: " + uploadServer)
	err := http.ListenAndServe(uploadServer, router)
	if err != nil {
		log.Println("upload server listen faild: " +err.Error())
		os.Exit(0)
	}
}

// 下载 server
func downloadServer()  {

	downloadServer := conf.GetString("listen.download")
	router := httprouter.New()
	httpHandle := NewHttpHandle()
	router.GET("/image/:name", httpHandle.ImageFind)

	log.Println("download server start listen: " + downloadServer)
	err := http.ListenAndServe(downloadServer, router)
	if err != nil {
		log.Println("download server listen faild:" + err.Error())
		os.Exit(0)
	}
}