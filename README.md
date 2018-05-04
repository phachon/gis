```
             _
   ____ _   (_)  _____
  / __  /  / /  / ___/
 / /_/ /  / /  (__  )
 \__, /  /_/  /____/
/____/

```
# go image server
go 实现的图片服务, 提供上传, 存储, 自动裁剪, 下载等功能

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/phachon/gis/)
[![license](https://img.shields.io/github/license/phachon/gis.svg?style=plastic)]()
[![download_count](https://img.shields.io/github/downloads/phachon/gis/total.svg?style=plastic)](https://github.com/phachon/gis/releases)
[![release](https://img.shields.io/github/release/phachon/gis.svg?style=plastic)](https://github.com/phachon/gis/releases)

## 功能
- http 上传
- 图片存储
- 按比例裁剪图片
- 图片下载浏览

## 安装

下载最新版本的二进制程序，下载地址：https://github.com/phachon/gis/releases

## 使用
- windows

```
gis.exe
# 指定配置文件启动
gis.exe --conf config.toml
```

- linux

```
./gis
# 指定配置文件启动
./gis --conf config.toml
```

## 配置

#### config.toml
```
[listen]
# 监听上传 server
upload="127.0.0.1:8087"
# 监听下载 server
download="127.0.0.1:8088"

[upload]
form_field="upload" // 表单提交字段
allow_type = [".jpg", ".jpeg", ".png"] // 允许上传的图片格式
max_size = 2048 // 图片的最大上传大小 KB
root_dir = "upload" // 图片上传根目录
filename_len = 16 // 图片保存文件名字符串长度
dirname_len = 4  // 目录树的目录名长度
thumbnails = ["200_200", "300_300", "200_400"] // 要生成的缩略图裁剪尺寸 width_height

[download]
# 下载的地址 协议://域名:端口
uri = "http://test.com:8088"

[appname] // appname 用于授权,可多个，app_key 需要和 客户端上传的 token 保持一致
    [appname.test]
    app_key = "ad%4a*a&ada@#ada"
    [appname.test1]
    app_key = "sd(4a*yu&dai#9d3"
```

## 接口说明

### 上传图片接口

- 请求地址：/image/upload?
- 请求方式：POST
- 请求 Header: Appname, Token (用来验证上传合法性)
- 返回格式：json

```
{
  "code": "1",   // 1:success, 0:error
  "message": "", // error message
  "data": {
       "image": "http://test.com:8088/image/LYEDBYKAFGGRJUFL.png"
       "image_200_200": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_200_200.png"
       "image_200_400": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_200_400.png"
       "image_300_300": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_300_300.png"
   }, // server image url
}
```

### Token 生成规则
```
token = md5(appname+appKey)
```

### 下载图片接口
- 接口地址：/image/:imageName
- 请求方式：GET
- 返回：图片

## 客户端调用示例
- [php](https://github.com/phachon/gis/tree/master/_example/php/upload.php)
- [go](https://github.com/phachon/gis/tree/master/_example/go/upload.go)
- [html](https://github.com/phachon/gis/tree/master/_example/html/index.html)

## 反馈

欢迎提交意见和代码，联系方式 phachon@163.com

## License

MIT

Thanks
---------
Create By phachon@163.com
