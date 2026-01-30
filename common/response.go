package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. 定义标准 JSON 结构
type Response struct {
	Code int         `json:"code"` // 200 表示成功
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

// 2. 成功时的封装
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// common/response.go

func Error(c *gin.Context, code int, msg interface{}) {
	var message string

	// 使用类型分支（Type Switch）来判断传入的是什么
	switch v := msg.(type) {
	case error:
		message = v.Error() // 如果是 error 接口，自动转字符串
	case string:
		message = v // 如果本来就是字符串，直接用
	default:
		message = "未知错误"
	}

	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}
