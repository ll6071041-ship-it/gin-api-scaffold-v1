package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/logger"
	"gin-api-scaffold-v1/middleware"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 1. 创建 Gin 实例 (使用 gin.New() 也就是一张白纸)
	r := gin.New()

	// 2. 注册全局中间件
	// ⚡️ 替换默认 logger 和 recovery
	r.Use(middleware.GinLogger(logger.Logger))
	r.Use(middleware.GinRecovery(logger.Logger, true))
	// ⚡️ 注册跨域中间件 (直接调用刚才写的函数)
	r.Use(middleware.Cors())

	// 3. 注册路由
	// 基础健康检查 (Ping)，通常用于 k8s 探针或负载均衡检测
	r.GET("/ping", controller.Ping)

	// 4. 业务路由分组 (标准做法)
	// 以后你的业务接口都放在 /api/v1 下面
	api := r.Group("/api/v1")
	{
		// 比如：api.GET("/users", controller.GetUser)
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from v1"})
		})
	}

	// 5. 处理 404 (当访问不存在的路径时)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "404 Not Found",
		})
	})

	return r
}
