package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-api-scaffold-v1/common" // 👈 引入我们封装好的 common 包
	"gin-api-scaffold-v1/dao"    // 引入 DAO 以便判断特定错误
	"gin-api-scaffold-v1/logic"
	"gin-api-scaffold-v1/models"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamSignUp

	// ShouldBindJSON 内部会进行两个动作：
	// A. 读取 JSON 绑定到结构体
	// B. 根据 tag (binding:"required") 进行校验
	if err := c.ShouldBindJSON(&p); err != nil {
		// 记录日志：这是开发看的，记录原始错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		// ⚡️ 核心改造：使用 common.Error
		// 我们把原始的 err 传进去，common.Error 内部会自动识别：
		// 如果是 validator 校验错误 -> 自动翻译成中文 (如 "密码必须大于6位")
		// 如果是 JSON 格式错误 -> 返回原始错误信息
		common.Error(c, common.CodeInvalidParam, err)
		return
	}

	// 2. 业务处理：调用 Logic 层
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))

		// ⚡️ 进阶处理：根据不同的错误类型，返回不同的业务状态码
		// 假设我们在 DAO 层定义了 var ErrorUserExist = errors.New("用户已存在")
		// 这里可以用 errors.Is 来判断
		if errors.Is(err, dao.ErrorUserExist) {
			common.Error(c, common.CodeUserExist, err)
			return
		}

		// 如果是其他未知错误（比如数据库挂了），就返回 "服务繁忙"
		common.Error(c, common.CodeServerBusy, err)
		return
	}

	// 3. 返回响应
	// 注册成功，不需要返回什么数据，传 nil 即可
	common.Success(c, nil)
}

// LoginHandler 处理登录请求的函数
func LoginHandler(c *gin.Context) {
	// 1. 获取参数
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 参数校验错误
		common.Error(c, common.CodeInvalidParam, err)
		return
	}

	// 2. 业务处理
	token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))

		// 登录失败通常有两种情况：用户不存在、密码错误
		// 为了安全，通常统称为 "用户名或密码错误" (CodeInvalidPassword)
		// 或者是根据 err 具体内容判断
		if err.Error() == "用户不存在" {
			common.Error(c, common.CodeUserNotExist, err)
		} else {
			common.Error(c, common.CodeInvalidPassword, err)
		}
		return
	}

	// 3. 返回响应
	// 将 Token 放在 Data 字段里返回给前端
	common.Success(c, gin.H{
		"token":   token,
		"user_id": 123456, // 举例：你也可以顺便把 userID 返回去
		"name":    p.Username,
	})
}

// Ping 心跳检测
func Ping(c *gin.Context) {
	// Ping 接口一般不需要复杂的结构，简单返回即可
	// 当然你也可以用 common.Success(c, "pong")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
