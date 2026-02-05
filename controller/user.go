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
// @Summary      用户注册
// @Description  处理用户注册请求，需要提供用户名、密码和确认密码
// @Tags         用户相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        object body  models.ParamSignUp  true  "注册参数"
// @Success      200  {object} common.Response "注册成功"
// @Router       /signup [post]
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
// @Summary      用户登录
// @Description  处理用户登录请求，返回JWT Token
// @Tags         用户相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        object body  models.ParamLogin  true  "登录参数"
// @Success      200  {object} common.Response{data=models.ParamLogin} "登录成功"
// @Router       /login [post]
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
// @Summary      获取个人信息
// @Description  获取当前登录用户的详细信息 (需要携带Token)
// @Tags         用户相关接口
// @Accept       application/json
// @Produce      application/json
// @Security     ApiKeyAuth
// @Success      200  {object} common.Response "成功返回用户信息"
// @Router       /home [get]
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
// @Summary      心跳检测
// @Description  检查服务是否正常运行
// @Tags         基础接口
// @Success      200  {object} map[string]string "pong"
// @Router       /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
