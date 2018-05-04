package controllers

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"strconv"
	"gis/app"
	"gis/app/utils"
	"errors"
)

type BaseController struct {

}

type JsonResult struct {
	Code int `json:"code"`
	Message interface{} `json:"message"`
	Data interface{} `json:"data"`
}

// get request content text string
func (baseService *BaseController) GetCtxString(ctx *fasthttp.RequestCtx, key string) string {
	return string(ctx.QueryArgs().Peek(key))
}

// get request content text bool
func (baseService *BaseController) GetCtxBool(ctx *fasthttp.RequestCtx, key string) bool {
	str := string(ctx.QueryArgs().Peek(key))
	if str == "1" {
		return true
	}
	return false
}

// get request content text int
func (baseService *BaseController) GetCtxInt(ctx *fasthttp.RequestCtx, key string) int {
	str := string(ctx.QueryArgs().Peek(key))
	i, _ := strconv.Atoi(str)
	return i
}

// get request content text float64
func (baseService *BaseController) GetCtxFloat64(ctx *fasthttp.RequestCtx, key string) float64 {
	str := string(ctx.QueryArgs().Peek(key))
	i, _ := strconv.Atoi(str)
	return float64(i)
}

// return json error
func (baseService *BaseController) jsonError(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 0, message, data)
}

// return json success
func (baseService *BaseController) jsonSuccess(ctx *fasthttp.RequestCtx, message interface{}, data interface{}) {
	baseService.jsonResult(ctx, 1, message, data)
}

// return json result
func (baseService *BaseController) jsonResult(ctx *fasthttp.RequestCtx, code int, message interface{}, data interface{}) {
	if message == nil {
		message = ""
	}
	if data == nil {
		data = map[string]string{}
	}

	res := JsonResult {
		Code:    code,
		Message: message,
		Data:    data,
	}

	jsonByte, err := json.Marshal(res)
	if err != nil {
		ctx.Write([]byte(err.Error()))
	} else {
		ctx.Write(jsonByte)
	}
}

// 验证 token 合法
// params: ctx fasthttp.RequestCtx
// return: error. auth success err is nil
func (baseService *BaseController) authToken(ctx *fasthttp.RequestCtx) (err error) {

	headerAppName := string(ctx.Request.Header.Peek("Appname"))
	headerToken := string(ctx.Request.Header.Peek("Token"))

	if headerAppName == "" {
		return errors.New("upload auth failed, appname error! ")
	}
	if headerToken == "" {
		return errors.New("upload auth failed, appname error! ")
	}

	appConf := app.Conf.GetStringMapString("appname."+headerAppName)
	if len(appConf) == 0 {
		return errors.New( "upload auth failed, appname error! ")
	}
	appKey, ok := appConf["app_key"]
	if !ok {
		return errors.New("upload auth failed, appname conf error! ")
	}

	token := utils.Md5Encode(headerAppName+appKey)
	if token != headerToken {
		return errors.New( "upload auth failed, token error! ")
	}

	return nil
}