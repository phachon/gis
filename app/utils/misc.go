package utils

import (
	"math/rand"
	"time"
	"crypto/md5"
	"encoding/hex"
)

//生成随机字符串
func GetRandomString(strLen int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//字符串分割生成目录
func StringToPath(str string, n int) string {
	strLen := len(str)
	if n >= strLen {
		return "/" + str
	}
	r := strLen % n
	path := ""
	for i:=0; i < strLen - r; i+=n {
		path = path + "/" + str[i:i+n]
	}
	if r != 0 {
		path = path + "/" + str[strLen-r:]
	}
	return path
}

// md5加密
func Md5Encode(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}