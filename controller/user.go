package controller

import (
	"gin-api-scaffold-v1/logic"
	"gin-api-scaffold-v1/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamSignUp
	// ShouldBindJSON 会自动根据 tag (binding:"required"等) 进行校验
	if err := c.ShouldBindJSON(&p); err != nil {
		// 校验失败，记录日志，并返回错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		// 假如你有 common 包，可以用 common.Error(c, 400, err)
		// 这里我用标准写法，确保你复制过去能直接用：
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 业务处理：调用 Logic 层
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		c.JSON(200, gin.H{
			"code": 1002,
			"msg":  "注册失败: " + err.Error(),
		})
		return
	}

	// 3. 返回响应
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "注册成功",
	})
}

// LoginHandler 处理登录请求的函数
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		c.JSON(200, gin.H{
			"code": 1001,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 业务处理：调用 Logic 层
	token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(200, gin.H{
			"code": 1002,
			"msg":  "登录失败: " + err.Error(),
		})
		return
	}

	// 3. 返回响应 (把 Token 给前端)
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// Ping 心跳检测 (保留你原来的)
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
