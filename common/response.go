package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	// 引入你的 validator 包
	myValidator "gin-api-scaffold-v1/pkg/validator"
)

// Response 定义标准 JSON 结构
type Response struct {
	Code ResCode     `json:"code"` // 引用 code.go 里的 ResCode
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// Error 错误返回
func Error(c *gin.Context, code ResCode, err error) {
	var response Response
	response.Code = code

	if err == nil {
		response.Msg = code.Msg()
		response.Data = nil
		c.JSON(http.StatusOK, response)
		return
	}

	// 判断是否为 Validator 校验错误
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		response.Code = CodeInvalidParam
		response.Msg = CodeInvalidParam.Msg()
		translations := errs.Translate(myValidator.Trans)
		response.Data = myValidator.RemoveTopStruct(translations)
	} else {
		// 普通错误
		response.Msg = code.Msg()
		// 如果你想调试时看具体错误，可以取消下面这行的注释
		// response.Data = err.Error()
	}

	c.JSON(http.StatusOK, response)
}

// ErrorWithMsg 自定义错误信息返回
func ErrorWithMsg(c *gin.Context, code ResCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
