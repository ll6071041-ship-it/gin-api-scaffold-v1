package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收 gin 框架默认的日志，用 zap 替代
// 作用：记录每一个请求的详细信息（路径、IP、耗时、状态码等）
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 开始计时
		start := time.Now()

		// 2. 获取请求路径 (例如: /api/v1/login)
		path := c.Request.URL.Path
		// 获取请求参数 (例如: ?username=admin)
		query := c.Request.URL.RawQuery

		// 3. ⚡️ 让请求继续往下走！
		// 去执行后续的中间件，或者去执行你的 controller 业务逻辑
		c.Next()

		// ==============================
		// 4. 业务处理完了，回来计算耗时
		// ==============================
		cost := time.Since(start)

		// 5. 收集并记录日志
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),                                 // 状态码 (200, 404, 500)
			zap.String("method", c.Request.Method),                               // 请求方法 (GET, POST, DELETE)
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // 请求参数
			zap.String("ip", c.ClientIP()),                                       // 客户端 IP
			zap.String("user-agent", c.Request.UserAgent()),                      // 浏览器标识 (Chrome/Edge/Postman)
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // Gin 内部捕获的错误
			zap.Duration("cost", cost),                                           // ⚡️ 核心指标：耗时
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
// 作用：防崩！如果程序哪里写错了导致崩溃，它能兜底，并记录错误堆栈。
// 参数 stack: 是否记录堆栈信息 (true: 记录详细报错位置; false: 只记录报错信息)
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// recover() 必须在 defer 里调用，用来捕获 panic
			if err := recover(); err != nil {
				// --------------------------------------------------------
				// 1. 判断是否是 "Broken Pipe" 错误
				// 这种错误通常是用户网络不好、突然关闭浏览器导致的，不是你的代码 bug
				// --------------------------------------------------------
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取完整的 HTTP 请求内容（方便你看是哪个请求搞崩了系统）
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 如果是 Broken Pipe，只记录简单日志，不打印堆栈，不需要改状态码
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // 记录错误到 Gin 上下文
					c.Abort()            // 终止后续操作
					return
				}

				// --------------------------------------------------------
				// 2. 处理真正的代码崩溃 (Panic)
				// --------------------------------------------------------
				if stack {
					// stack = true 时：打印详细堆栈 (debug.Stack())
					// 这会告诉你具体是 main.go 第几行代码出错了
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 👈 核心：报错的“案发现场”
					)
				} else {
					// stack = false 时：只告诉你有错，不告诉你错哪了
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 3. 返回 500 给前端，表示服务器内部错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 继续执行请求
		c.Next()
	}
}
