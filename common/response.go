package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	// 👇 引入我们需要用到的自定义 validator 包 (里面有 RemoveTopStruct 和 Trans)
	myValidator "gin-api-scaffold-v1/pkg/validator"
)

// 1. 定义标准 JSON 结构
type Response struct {
	Code int         `json:"code"` // 业务状态码 (200=成功, 400=参数错误, 500=系统错误)
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据 (可能是对象、列表，或者错误详情 map)
}

// 2. 成功时的封装
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// 3. 错误处理封装 (🔥 核心改造部分)
// c: 上下文
// code: 业务错误码 (比如 1001)
// err: 具体的错误对象
func Error(c *gin.Context, code int, err error) {
	var response Response
	response.Code = code

	// =========================================================
	// 🔥 关键点：类型断言 (Type Assertion)
	// 我们判断传入的 err 到底是不是 "参数校验错误" (validator.ValidationErrors)
	// =========================================================
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// Case A: 如果不是校验错误 (比如数据库连不上、逻辑错误)
		// 直接返回错误的字符串描述
		response.Msg = err.Error()
		response.Data = nil
	} else {
		// Case B: 如果是参数校验错误！

		// 1. 使用我们在 pkg/validator 里初始化的全局翻译器 Trans 进行翻译
		//    这会返回一个 map[string]string，key是字段名，value是中文错误
		translations := errs.Translate(myValidator.Trans)

		// 2. 去除结构体名字前缀
		//    把 "SignUpParam.Age" 变成 "age"
		cleanData := myValidator.RemoveTopStruct(translations)

		// 3. 构造返回
		//    Msg 提示通用信息 "请求参数错误"
		//    Data 里放具体的字段错误详情，方便前端展示在输入框下面
		response.Msg = "请求参数错误"
		response.Data = cleanData
	}

	c.JSON(http.StatusOK, response)
}
