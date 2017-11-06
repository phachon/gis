package main

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
)

// go 实现的图片上传,存储,缩放,下载服务

func init()  {
	initConfig();
}

func main()  {

	server := conf.GetString("listen.server")
	router := httprouter.New()
	httpHandle := NewHttpHandle()
	router.GET("/", httpHandle.Index)
	router.POST("/image/upload", httpHandle.ImageUpload)
	router.GET("/image/:name", httpHandle.ImageFind)

	log.Println("start server: " + server)
	err := http.ListenAndServe(server, router)
	if err != nil {
		log.Println(err.Error())
	}
}