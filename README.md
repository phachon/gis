# go-imageServer
go 实现的图片服务, 提供基本的上传,存储,缩放,下载等功能

## 功能
- http 上传
- 图片存储
- 按比例自动缩放生成图片
- 图片浏览

## 安装
1. 普通安装<br>
下载地址：https://github.com/phachon/go-imageServer/releases<br>
找到对应的版本下载
2. 手动安装<br>
http下载地址：https://github.com/phachon/go-imageServer.git<br>
ssh 下载地址：git@github.com:phachon/go-imageServer.git<br>
```
cd go-imageServer
go get ./...
go build ./
```
## 访问
- windows:<br>
go-imageServer.exe
- linux:<br>
./go-imageServer

## 配置文件

config.toml 文件读取顺序：<br>
- /etc/go-imageServer/config.toml
- $HOME/.go-imageServer/config.toml
- ./config.toml

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
thumbnails = ["200_200", "300_300", "200_400"] // 要生成的缩略图尺寸 width_height

[download]
# 下载的地址 协议://域名:端口
uri = "http://test.com:8088"

[appname] // appname 用于授权,可多个
    [appname.test]
    app_key = "ad%4a*a&ada@#ada"
    [appname.test1]
    app_key = "sd(4a*yu&dai#9d3"
```

## api 说明

- 上传接口
地址：/image/upload?<br>
请求方式：POST<br>
请求 header ：Appname, Token (用来验证上传合法性)<br>
返回：json<br>
```
{
  "code":"1",   // 1:success, 0:error
  "message":"", // error message
  "data":{
       "image": "http://test.com:8088/image/LYEDBYKAFGGRJUFL.png"
       "image_200_200": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_200_200.png"
       "image_200_400": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_200_400.png"
       "image_300_300": "http://test.com:8088/image/LYEDBYKAFGGRJUFL_300_300.png"
   }, // server image url
}
```

- 访问接口
地址：/image/:imageName<br>
请求方式：GET<br>
返回：图片

- Token 生成规则
```
token = md5(appname+appKey)
```

## 客户端调用示例
- php <br>
https://github.com/phachon/go-imageServer/example/php/php.go
- go <br>
https://github.com/phachon/go-imageServer/example/go/upload.go
- html <br>
https://github.com/phachon/go-imageServer/example/html/index.html


## 反馈

欢迎提交意见和代码，联系方式 phachon@163.com

## License

MIT

Thanks
---------
Create By phachon@163.com
