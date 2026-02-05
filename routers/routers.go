package routers

import (
	"net/http"
	"time" // 👈 【新增】需要用到时间计算

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper" // 👈 【新增】需要读取配置文件

	// 👇 【新增】这里必须导入 swagger 的两个包，否则下面的 gs 和 swaggerFiles 会报错 undefined
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	"gin-api-scaffold-v1/controller"
	"gin-api-scaffold-v1/middleware"
)

// SetupRouter 配置路由入口
// 负责把所有的 URL 路径和 Controller 里的函数对应起来
func SetupRouter() *gin.Engine {
	// 1. 创建 Gin 实例 (白纸模式)
	// 使用 gin.New() 而不是 gin.Default()，以便我们自己定制中间件
	r := gin.New()

	// =======================================================
	// 2. 注册全局中间件 (Middleware)
	// =======================================================
	// 记录请求日志：把 Gin 的请求详情记录到我们的 Zap 日志文件中
	r.Use(middleware.GinLogger())
	// 崩溃恢复：防止程序 Panic 导致整个服务挂掉
	r.Use(middleware.GinRecovery(true))
	// 跨域处理 (CORS)：允许前端跨域访问
	r.Use(middleware.Cors())

	// 🔥 【新增】注册全局限流中间件 (令牌桶)
	// 从配置文件读取 QPS (每秒请求数)
	qps := viper.GetInt64("rate_limit.qps")
	if qps > 0 {
		// 计算填充间隔: 如果 QPS 是 1000，那么间隔就是 1秒/1000 = 1毫秒
		// 也就是说：每 1 毫秒往桶里放一个令牌，一秒钟正好放 1000 个
		fillInterval := time.Second / time.Duration(qps)

		// 容量也设置为 QPS 的大小，允许瞬间爆发 1000 个请求
		r.Use(middleware.RateLimitMiddleware(fillInterval, qps))
	}

	// =======================================================
	// 3. 注册基础路由 (Infrastructure)
	// =======================================================
	// 健康检查接口，访问：GET /ping
	r.GET("/ping", controller.Ping)

	// =======================================================
	// 4. 业务路由分组 (Business Logic)
	// =======================================================
	// 创建一个路由组，前缀是 /api/v1
	// 此时 api 变量还没有挂载 JWT 中间件
	api := r.Group("/api/v1")
	{
		// ---------------------------------------------------
		// 🚫 公开路由 (无需 Token 即可访问)
		// ---------------------------------------------------
		// 用户注册：POST /api/v1/signup
		api.POST("/signup", controller.SignUpHandler)
		// 用户登录：POST /api/v1/login
		api.POST("/login", controller.LoginHandler)

		// ---------------------------------------------------
		// 🔒 私有路由 (必须带 Token 才能访问)
		// ---------------------------------------------------
		// ⚡️ 核心技巧：创建一个新的路由组 auth
		// 虽然 auth 的路径前缀和 api 一样 (都是 /api/v1)，
		// 但我们只给 auth 这个组挂载了 JWT 中间件！
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware()) // 挂载鉴权中间件
		{
			// 获取个人信息 (测试 JWT 用)
			// 访问路径：GET /api/v1/home
			// 只有 Token 验证通过，才会进入 controller.GetProfileHandler
			auth.GET("/home", controller.GetProfileHandler)

			// 未来其他的私有接口写在这里...
			// auth.POST("/article/publish", controller.CreateArticleHandler)
		}
	}

	// =======================================================
	// 5. 处理 404 (Not Found)
	// =======================================================
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "404 Not Found (你访问的路径不存在)",
		})
	})

	// =======================================================
	// 6. 注册 Swagger 文档路由
	// =======================================================
	// 访问地址：http://localhost:port/swagger/index.html
	// gs 和 swaggerFiles 现在可以正常使用了，因为我们在文件顶部 import 了它们
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	return r
}
