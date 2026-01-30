package controller

import (
	"github.com/gin-gonic/gin"
)

// Ping 处理函数
// @Summary 只是一个测试接口
func Ping(c *gin.Context) {
	// 简单的返回 JSON，证明服务活着
	c.JSON(200, gin.H{
		"message": "pong",
		"status":  "success",
	})
}
