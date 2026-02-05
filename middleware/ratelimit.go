package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware 令牌桶限流中间件
// fillInterval: 往桶里放令牌的时间间隔 (比如 time.Second / 100 表示每秒放100个)
// capacity: 桶的容量 (Capacity)，即允许瞬间爆发的最大并发数
func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	// 创建一个令牌桶
	// 参数1: 填充间隔, 参数2: 容量
	bucket := ratelimit.NewBucket(fillInterval, capacity)

	return func(c *gin.Context) {
		// 尝试拿 1 个令牌
		// TakeAvailable(1) 是非阻塞的，如果桶里有令牌就返回 1，没有就返回 0
		if bucket.TakeAvailable(1) < 1 {
			// 拿不到令牌，直接拒绝
			c.JSON(http.StatusOK, gin.H{
				"code": 429, // 429 Too Many Requests
				"msg":  "请求太快了，服务器繁忙，请稍后再试",
			})
			c.Abort() // 🛑 拦截请求，不让它往后走了
			return
		}
		// 拿到了，放行
		c.Next()
	}
}
