package middleware

import (
	"strings"

	"gin-api-scaffold-v1/common"
	"gin-api-scaffold-v1/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization Header
		// 行业规范：前端要把 Token 放在 Header 的 "Authorization" 字段里
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 如果没带 Token，直接拒绝，返回 "需要登录"
			common.Error(c, common.CodeNeedLogin, nil)
			c.Abort() // 🚫 阻止执行后续函数
			return
		}

		// 2. 解析 Header 格式
		// 行业规范：Authorization: Bearer <token>
		// 所以我们要按空格切割，取第2部分
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			common.Error(c, common.CodeInvalidToken, nil)
			c.Abort()
			return
		}

		// 3. 解析 Token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			// Token 过期或无效
			common.Error(c, common.CodeInvalidToken, err)
			c.Abort()
			return
		}

		// 4. ✅ 验证通过！将当前请求的 UserID 信息保存到上下文 c 中
		// 这样后续的 Controller 就能知道是谁在访问了
		c.Set("userID", mc.UserID)
		c.Set("username", mc.Username)

		c.Next() // 放行，进入下一个环节
	}
}
