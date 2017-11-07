package main

import (
	"net/http"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"mime/multipart"
	"bytes"
	"log"
	"io"
)

func main()  {

	// 上传地址
	uploadUrl := "http://127.0.0.1:8087/image/upload"
	// 预分配的 appname
	appname := "test"
	// 预分配 appname 对应的 appKey
	appKey := "ad%4a*a&ada@#ada"
	// 加密算法
	token := md5Encode(appname+appKey)
	// 设置 header 头
	headers := map[string]string{
		"Appname": appname,
		"Token": token,
	}
	// 注意：根据自己的目录对应的图片的绝对路径
	file := "../image/test.jpg"

	fmt.Println(file)
	os.Exit(1)
	body, code, err := httpUpload(uploadUrl, file, headers)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println(code)
		fmt.Println(body)
	}
}

// http upload
func httpUpload(queryUrl string, file string, headerValues map[string]string) (body string, code int, err error) {

	// 创建表单文件, 第一个参数是字段名，第二个参数是文件名
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("upload", file)
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
		return
	}

	// 从文件读取数据，写入表单
	srcFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("Open source file failed: %s\n", err)
		return
	}
	defer srcFile.Close()

	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
		return
	}
	writer.Close()

	// 发送表单数据
	req, err := http.NewRequest("POST", queryUrl, buf)
	if err != nil {
		log.Fatalf("Request falied: %s\n", err)
		return
	}

	// 设置 http header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if (headerValues != nil) && (len(headerValues) > 0) {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}

	// Client 发送
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(bodyByte), code, nil
}

// md5加密
func md5Encode(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
