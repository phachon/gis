package main

import (
	"net/http"
	"log"
	"os"
	"strings"
	"io"
)

// go 实现的图片上传,存储,缩放,下载服务

//初始化配置
func init()  {
	initConfig();
}

func main()  {
	
	server := conf.GetString("listen.server")
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/upload/image", HandleImage)
	
	log.Println("start server: "+server)
	err := http.ListenAndServe(server, nil)
	if err != nil {
		log.Println(err.Error())
	}
	
}

func HandleImage(w http.ResponseWriter, req *http.Request) {
	
	//处理上传请求
	if req.Method != "POST" {
		return
	}
	
	formField := conf.GetString("upload.form_field")
	allowTypeSlice := conf.GetStringSlice("upload.allow_type")
	rootDir := conf.GetString("upload.root_dir")
	filenameLen := conf.GetInt("upload.filename_len")
	dirNameLen := conf.GetInt("upload.dirname_len")
	
	req.ParseMultipartForm(1024)
	file, fileHeader, err := req.FormFile(formField)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()
	
	//文件上传格式判断
	filename := fileHeader.Filename
	ext := filename[strings.LastIndex(filename, "."):]
	isAllow := false
	for _, allowType := range allowTypeSlice {
		if strings.ToLower(allowType) == strings.ToLower(ext) {
			isAllow = true
			break;
		}
	}
	if isAllow == false {
		log.Println("no allow upload type!")
		return
	}
	
	// 生成文件保存路径
	randString := strings.ToUpper(NewUtils().GetRandomString(filenameLen))
	uploadPath := rootDir + NewUtils().StringToPath(randString, dirNameLen)
	err = os.MkdirAll(uploadPath, 0755)
	if err != nil {
		log.Println(err.Error())
		return
	}
	saveFilename := uploadPath + "/" + randString + ext
	_, err = os.Stat(saveFilename)
	if(!os.IsNotExist(err)) {
		log.Println("file exist!")
		return
	}
	
	//将文件写入到指定的位置
	f, err := os.OpenFile(saveFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	
	width := 200;
	height := 200;
	thumbSaveFilename := uploadPath + "/" + randString + "_200_200" +ext
	//开始生成缩略图
	NewImage().Scaling(saveFilename, thumbSaveFilename, width, height)
	
	w.Write([]byte("upload ok!"))
}