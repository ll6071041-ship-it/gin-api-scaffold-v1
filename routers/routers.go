package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/middleware"
)

// SetupRouter 配置路由入口
// 负责把所有的 URL 路径和 Controller 里的函数对应起来
func SetupRouter() *gin.Engine {
	// =======================================================
	// 1. 创建 Gin 实例 (白纸模式)
	// =======================================================
	// 使用 gin.New() 而不是 gin.Default()。
	// 默认的 gin.Default() 强行绑定了标准库的 Logger，
	// 我们已经有了更高级的 Zap Logger，所以要用 gin.New() 从零开始配置。
	r := gin.New()

	// =======================================================
	// 2. 注册全局中间件 (Middleware)
	// =======================================================

	// 记录请求日志：把 Gin 的请求详情记录到我们的 Zap 日志文件中
	r.Use(middleware.GinLogger())

	// 崩溃恢复：如果程序哪里写得不对导致 Panic (崩溃)，
	// 这个中间件会兜底捕获，防止整个服务器挂掉，并打印错误堆栈。
	r.Use(middleware.GinRecovery(true))

	// 跨域处理 (CORS)：
	// 允许前端 (Vue/React) 跨域访问我们的 API，防止出现 "CORS error"
	r.Use(middleware.Cors())

	// =======================================================
	// 3. 注册基础路由 (Infrastructure)
	// =======================================================

	// 健康检查接口
	// 访问：GET /ping
	// 作用：通常用于 K8s 探针、负载均衡器检测服务是否存活
	r.GET("/ping", controller.Ping)

	// =======================================================
	// 4. 业务路由分组 (Business Logic)
	// =======================================================
	// 建议所有业务接口都放在 /api/v1 下面，方便未来升级 v2 版本
	api := r.Group("/api/v1")
	{
		// ---------------------------------------------------
		// 用户业务模块 (User Module)
		// ---------------------------------------------------

		// 用户注册
		// 访问：POST /api/v1/signup
		// 逻辑：校验参数 -> 查重 -> 生成ID -> 密码加密 -> 存库
		api.POST("/signup", controller.SignUpHandler)

		// 用户登录
		// 访问：POST /api/v1/login
		// 逻辑：校验参数 -> 查库 -> 密码比对 -> 颁发 Token
		api.POST("/login", controller.LoginHandler)

		// ---------------------------------------------------
		// 如果以后有其他模块 (比如订单、商品)，继续在这里往下写
		// ---------------------------------------------------
		// api.GET("/products", controller.GetProductList)
		// api.POST("/orders", controller.CreateOrder)
	}

	// =======================================================
	// 5. 处理 404 (Not Found)
	// =======================================================
	// 当用户访问了一个不存在的路径时，返回友好的 JSON 提示
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "404 Not Found (你访问的路径不存在)",
		})
	})

	return r
}
