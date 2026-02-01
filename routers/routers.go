package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/middleware"
)

// SetupRouter 配置路由入口
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.GinLogger())
	r.Use(middleware.GinRecovery(true))
	r.Use(middleware.Cors())

	// 基础路由
	r.GET("/ping", controller.Ping)

	// 业务路由分组 /api/v1
	api := r.Group("/api/v1")
	{
		// =======================================================
		// 🚫 公开路由 (无需 Token 即可访问)
		// =======================================================

		api.POST("/signup", controller.SignUpHandler) // 注册
		api.POST("/login", controller.LoginHandler)   // 登录

		// =======================================================
		// 🔒 私有路由 (必须带 Token 才能访问)
		// =======================================================

		// ⚡️ 创建一个新的路由组，专门挂载 JWT 中间件
		// 只有进入这个组的请求，才会被 JWTAuthMiddleware 拦截检查
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			// 举例：获取首页/个人信息
			// 只有 Token 验证通过，才会执行里面的逻辑
			auth.GET("/home", func(c *gin.Context) {
				// 从上下文中取出中间件塞进去的 userID 和 username
				userID, _ := c.Get("userID")
				username, _ := c.Get("username")

				c.JSON(200, gin.H{
					"code": 1000,
					"msg":  "success",
					"data": gin.H{
						"id":   userID,
						"name": username,
						"info": "你能看到这条信息，说明你已经登录成功了！",
					},
				})
			})

			// 以后其他的需要登录的接口都写在这里
			// auth.POST("/article", controller.CreateArticle)
		}
	}

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "404 Not Found",
		})
	})

	return r
}
