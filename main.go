package main

import (
	"go-imageServer/app"
	"go-imageServer/app/controllers"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

// gis (go image server)


var (
	imageController = controllers.NewImageController()
)

func main()  {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				app.Log.Errorf("gis upload server crash, %v", e)
			}
		}()
		uploadServer()
	}()

	downloadServer()
}

// 上传 server
func uploadServer()  {

	uploadServer := app.Conf.GetString("listen.upload")
	app.Log.Info("start listen server "+uploadServer)

	router := fasthttprouter.New()
	router.POST("/image/upload", imageController.Upload)
	router.OPTIONS("/image/upload", imageController.CrossDomain)

	app.Log.Infof("upload server start listen: %s", uploadServer)

	err := fasthttp.ListenAndServe(uploadServer, router.Handler)
	if err != nil {
		app.Log.Errorf("listen upload server %s error: %s", uploadServer, err.Error())
	}
}

// 下载 server
func downloadServer()  {

	downloadServer := app.Conf.GetString("listen.download")
	router := fasthttprouter.New()

	router.GET("/image/:name", imageController.Download)
	router.OPTIONS("/image/:name", imageController.CrossDomain)

	app.Log.Infof("download server start listen: %s", downloadServer)

	err := fasthttp.ListenAndServe(downloadServer, router.Handler)
	if err != nil {
		app.Log.Info("listen download server "+downloadServer+" error :"+err.Error())
	}
}