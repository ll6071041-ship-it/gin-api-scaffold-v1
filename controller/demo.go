package controller

import (
	"fmt"
	"gin-api-scaffold-v1/common"
	"gin-api-scaffold-v1/models" // 👈 1. 引入 models 包

	"github.com/gin-gonic/gin"
)

// TestValidator 测试参数校验功能的接口
func TestValidator(c *gin.Context) {
	// 👈 2. 使用 models.RegisterParam
	var p models.RegisterParam

	// 1. ShouldBindJSON 会根据 Content-Type 读取 JSON 并绑定到结构体
	//    同时会根据 tag (binding:"...") 进行校验
	if err := c.ShouldBindJSON(&p); err != nil {
		// 2. 校验失败！
		//    直接把原始的 err 丢给 common.Error
		common.Error(c, 400, err)
		return
	}

	// 3. 校验通过，处理业务逻辑
	fmt.Printf("注册成功: %+v\n", p)

	// 注意：如果你想返回 p 里面的字段，也需要用 p.Username 这样访问
	common.Success(c, gin.H{
		"user_id": 12345,
		"name":    p.Username,
	})
}

// 原来的 Ping 函数保留
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"status":  "success",
	})
}
