package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response
// @Description: 响应结构体
type Response struct {
	Suc  bool        `json:"suc"`  // 是否成功
	Code uint        `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
}

func RespOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Suc:  true,
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func RespFail(ctx *gin.Context, code uint, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Suc:  false,
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func Sse(c *gin.Context, event string, data interface{}) {
	if event == "last" {
		c.SSEvent("message", Response{
			Suc:  true,
			Code: 1,
			Msg:  "ok",
			Data: data,
		})
		return
	}
	if event == "message" {
		c.SSEvent("message", Response{
			Suc:  true,
			Code: 0,
			Msg:  "ok",
			Data: data,
		})
	}
	if event == "error" {
		c.SSEvent("error", Response{
			Suc:  false,
			Code: 0,
			Msg:  "error",
			Data: data,
		})
	}
}
