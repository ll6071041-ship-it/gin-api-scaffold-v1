package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-api-scaffold-v1/common"
	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/logic"
	"gin-api-scaffold-v1/models"
)

// SignUpHandler 处理注册请求
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		common.Error(c, common.CodeInvalidParam, err)
		return
	}

	// 2. 业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, dao.ErrorUserExist) {
			common.Error(c, common.CodeUserExist, err)
			return
		}
		common.Error(c, common.CodeServerBusy, err)
		return
	}

	// 3. 返回响应
	common.Success(c, nil)
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	// 1. 获取参数
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		common.Error(c, common.CodeInvalidParam, err)
		return
	}

	// 2. 业务处理
	token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if err.Error() == "用户不存在" {
			common.Error(c, common.CodeUserNotExist, err)
		} else {
			common.Error(c, common.CodeInvalidPassword, err)
		}
		return
	}

	// 3. 返回响应
	// ⚡️ 这里我们返回 Token 和用户名
	// 前端拿到 Token 后，会自动解码出 UserID，所以这里不传 UserID 也可以
	common.Success(c, gin.H{
		"token":    token,
		"username": p.Username,
	})
}

// GetProfileHandler 获取用户个人信息 (测试 JWT 用)
// ⚡️ 这是新增的，用来替代 routers.go 里那个匿名函数
func GetProfileHandler(c *gin.Context) {
	// 1. 从上下文中取出 userID (这是中间件 middleware.JWTAuthMiddleware 塞进去的)
	// 如果取不到，说明中间件没生效（或者没配置好），属于系统级错误
	userID, exists := c.Get("userID")
	if !exists {
		zap.L().Error("GetProfileHandler: userID not found in context")
		common.Error(c, common.CodeNeedLogin, nil)
		return
	}

	// 2. 取出 username
	username, _ := c.Get("username")

	// 3. (可选) 这里通常会拿 userID 去数据库查更详细的信息
	// user, err := logic.GetUserByID(userID.(int64))

	// 4. 返回数据
	common.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
		"message":  fmt.Sprintf("你好 %s，Token 验证成功！这是你的私密数据。", username),
	})
}

// Ping 心跳检测
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
