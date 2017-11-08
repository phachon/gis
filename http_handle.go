package main

import (
	"github.com/julienschmidt/httprouter"
	"strings"
	"os"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"fmt"
	"io"
)

// 获取文件大小的接口
type Size interface {
	Size() int64
}
// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

type HttpHandle struct {

}

func NewHttpHandle() *HttpHandle  {
	return &HttpHandle{}
}

// 处理首页静态文件
func (handle *HttpHandle) Index(w http.ResponseWriter, req *http.Request, params httprouter.Params)  {
	//fs := http.FileServer(http.Dir("public"))
	//fs.ServeHTTP(w, req)
	return
}

// 验证 token 合法
// params: req http.request
// return: error. auth success err is nil
func (handle *HttpHandle) authToken(req *http.Request) (err error) {

	headerAppname := req.Header.Get("Appname")
	headerToken := req.Header.Get("Token")

	if headerAppname == "" {
		return fmt.Errorf("%s", "auth faild: appname error!")
	}
	if headerToken == "" {
		return fmt.Errorf("%s", "auth faild: appname error!")
	}

	appConf := conf.GetStringMapString("appname."+headerAppname)
	if len(appConf) == 0 {
		return fmt.Errorf("%s", "auth faild: appname error!")
	}
	appKey, ok := appConf["app_key"]
	if !ok {
		return fmt.Errorf("%s", "auth faild: appname conf error!")
	}

	token := NewUtils().Md5Encode(headerAppname+appKey)
	if token != headerToken {
		return fmt.Errorf("%s", "auth faild: token error!")
	}

	return nil
}

//处理跨域
func (handle *HttpHandle) CrossDomain(w http.ResponseWriter, req *http.Request, params httprouter.Params)  {

	if req.Method == "OPTIONS" {
		if origin := req.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token, Appname")
	}
	log.Println("allow crossdomain")
}

// 处理上传请求
func (handle *HttpHandle) ImageUpload(w http.ResponseWriter, req *http.Request, params httprouter.Params) {

	err := handle.authToken(req)
	if err != nil {
		log.Println(err.Error())
		handle.jsonError(w, err.Error(), nil)
		return
	}
	formField := conf.GetString("upload.form_field")
	allowTypeSlice := conf.GetStringSlice("upload.allow_type")
	rootDir := conf.GetString("upload.root_dir")
	filenameLen := conf.GetInt("upload.filename_len")
	dirNameLen := conf.GetInt("upload.dirname_len")
	maxSize := conf.GetInt("upload.max_size")
	thumbnails := conf.GetStringSlice("upload.thumbnails")
	server := conf.GetString("listen.server")
	imageUrl := "http://"+server+"/image/"

	req.ParseMultipartForm(4*1024)
	file, fileHeader, err := req.FormFile(formField)
	if err != nil {
		log.Println("upload field error: ", err.Error())
		handle.jsonError(w, "Upload field error!", nil)
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
		log.Println("Forbidden upload format: " + ext)
		handle.jsonError(w, "Forbidden upload format!", nil)
		return
	}

	// 判断上传文件大小
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		size := fileInfo.Size()/1024
		if size > int64(maxSize) {
			log.Printf("Upload image beyond maximum limit: %d kb", maxSize)
			handle.jsonError(w, "Upload image size"+ strconv.Itoa(int(size)) +"maximum limit! ", nil)
			return
		}
	}

	// 生成文件保存路径(防止发生随机碰撞)
	var i int
	var randString string
	var uploadPath string
	var saveFilename string
	for i = 0; i < 10; i++ {
		randString = strings.ToUpper(NewUtils().GetRandomString(filenameLen))
		uploadPath = rootDir + NewUtils().StringToPath(randString, dirNameLen)
		err = os.MkdirAll(uploadPath, 0755)
		if err != nil {
			log.Println("Create upload dir failed: " + err.Error())
			handle.jsonError(w, "Create upload dir failed!", nil)
			return
		}
		saveFilename = uploadPath + "/" + randString + ext
		_, err = os.Stat(saveFilename)
		if(os.IsNotExist(err)) {
			break
		}
	}
	if i == 10 {
		log.Println("Create upload dir failed: random collision 10!")
		handle.jsonError(w, "Create upload dir failed!", nil)
		return
	}

	//将文件写入到指定的位置
	f, err := os.OpenFile(saveFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Save file error: " + err.Error())
		handle.jsonError(w, "Save image file error!", nil)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	data := map[string]string{
		"image": imageUrl + randString + ext,
	}

	//开始生成缩略图
	for _, thumbnail := range thumbnails {
		thumbnailSlice := strings.Split(thumbnail, "_")
		if len(thumbnailSlice)  != 2 {
			continue
		}
		width, _:= strconv.Atoi(thumbnailSlice[0])
		height, _:= strconv.Atoi(thumbnailSlice[1])
		thumbSaveFilename := uploadPath + "/" + randString + "_" + thumbnail + ext
		err := NewImager().Scaling(saveFilename, thumbSaveFilename, width, height)
		if err != nil {
			log.Println("make thumbnail image error: " + err.Error())
			continue
		}

		data["image_"+thumbnail] = imageUrl+ randString + "_" + thumbnail + ext
	}

	handle.jsonSuccess(w, "", data)
}

// 处理查看请求
func (handle *HttpHandle) ImageFind(w http.ResponseWriter, req *http.Request, params httprouter.Params)  {

	name := params.ByName("name")
	rootDir := conf.GetString("upload.root_dir")
	dirNameLen := conf.GetInt("upload.dirname_len")

	if !strings.Contains(name, ".") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found!"))
		return
	}

	ext := name[strings.LastIndex(name, "."):]
	if ext == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found!"))
		return
	}
	filename := name[:strings.LastIndex(name, ".")]
	n := strings.Index(name, "_")

	if n != -1 {
		filename = filename[:n]
	}

	imagePath := rootDir + NewUtils().StringToPath(filename, dirNameLen)
	imagePath = imagePath + "/" + name

	http.ServeFile(w, req, imagePath)
}

// 返回 error json
func (handle *HttpHandle) jsonError(w http.ResponseWriter, message string, data interface{}) {
	handle.jsonMessage(w, 0, message, data)
}

// 返回 success json
func (handle *HttpHandle) jsonSuccess(w http.ResponseWriter, message string, data interface{})  {
	handle.jsonMessage(w, 1, message, data)
}

// 返回 json
func (handle *HttpHandle) jsonMessage(w http.ResponseWriter, code int, message, data interface{}) {
	w.Header().Set("content-type", "application/json")

	type Result struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
		Data    interface{} `json:"data"`
	}
	result := Result{
		Code:    code,
		Message: message,
		Data:    data,
	}
	resultByte, err := json.Marshal(result)
	if err != nil {
		result.Code = 0
		result.Message = err.Error()
		resultByte, _ = json.Marshal(result)
		w.Write(resultByte)
		return
	}
	//responseByte := callback+"("+string(resultByte)+")"
	w.Write(resultByte)
}