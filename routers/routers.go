package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/middleware"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 1. 创建 Gin 实例
	// 使用 gin.New() 而不是 gin.Default()
	// 原因：gin.Default() 会自动加载默认的 Logger 和 Recovery 中间件
	// 我们已经自己用 Zap 实现了这两个功能，所以需要一张“白纸” (gin.New())
	r := gin.New()

	// 2. 注册全局中间件
	// =======================================================
	// ⚡️ 核心修改：这里不再需要传入 logger.Logger 参数
	// ⚡️ 它们内部现在会自动使用全局的 zap.L()
	// =======================================================

	// 记录请求日志 (替代 Gin 默认的输出)
	r.Use(middleware.GinLogger())

	// 捕获 Panic 防止崩溃 (true 表示打印详细错误堆栈，方便排错)
	r.Use(middleware.GinRecovery(true))

	// 处理跨域请求 (让前端能正常调用接口)
	r.Use(middleware.Cors())

	// 3. 注册基础路由
	// 基础健康检查 (Ping)，通常用于 k8s 探针或负载均衡检测
	r.GET("/ping", controller.Ping)

	// 4. 业务路由分组 (API Versioning)
	// 建议所有业务接口都放在 /api/v1 下面，方便未来升级 v2 版本
	api := r.Group("/api/v1")
	{
		// 测试接口 (原有的)
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello from v1",
				"user_id": c.Query("user_id"), // 例子：获取参数
			})
		})

		// =====================================================================
		// 🔥 新增：参数校验测试接口
		// =====================================================================
		// 对应 controller/demo.go 中的 TestValidator 函数
		// 发送 POST 请求到 /api/v1/validator_test，Body 带上 JSON 数据即可测试
		api.POST("/validator_test", controller.TestValidator)

		// 可以在这里继续添加其他业务路由，例如：
		// api.POST("/login", controller.Login)
		// api.POST("/register", controller.Register)
	}

	// 5. 处理 404 (当访问不存在的路径时)
	// 这是一个好习惯，返回 JSON 格式的 404，而不是默认的纯文本
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "404 Not Found (没有找到该路径)",
		})
	})

	return r
}
