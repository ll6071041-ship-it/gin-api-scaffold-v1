package routers

import (
	"github.com/gin-gonic/gin"

	// 👇 注意：如果你现在的 go.mod 里的 module 还没改，就还是用 gin-api-scaffold-v1-v1
	// 如果你已经打算叫 gin-api-scaffold，这里记得改成 gin-api-scaffold-v1/controller
	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/logger"
	"gin-api-scaffold-v1/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 1. 基础中间件 (日志 + 恢复)
	r.Use(middleware.GinLogger(logger.Logger), middleware.GinRecovery(logger.Logger, true))

	// 2. 跨域配置 (保留)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // 加了 Authorization
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 3. 注册路由
	// ❌ 删掉原来的 Todo 路由
	// ✅ 只保留一个基础的 Ping 接口，证明脚手架能通
	r.GET("/ping", controller.Ping)

	// 如果你想保留 v1 分组的结构，也可以这样写：
	// v1 := r.Group("/v1")
	// {
	//     v1.GET("/ping", controller.Ping)
	// }

	return r
}
