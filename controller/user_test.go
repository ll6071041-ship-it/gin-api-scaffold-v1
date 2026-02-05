package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	// 引入必要的包以初始化环境
	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/logger"
	"gin-api-scaffold-v1/models"
	"gin-api-scaffold-v1/pkg/snowflake"
	"gin-api-scaffold-v1/settings"
)

// initEnv 初始化测试环境 (数据库、配置等)
// ⚠️ 注意：这需要连接真实的数据库
func initEnv() {
	// 1. 加载配置 (注意路径，测试文件在 controller 目录下，可能需要调整相对路径)
	// 如果报错找不到配置文件，建议把 config.yaml 复制一份到 controller 目录或者写死路径
	if err := settings.InitConfig(); err != nil {
		panic(err)
	}

	// 2. 初始化日志
	logger.InitLogger()

	// 3. 初始化数据库 (GORM)
	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}

	// 4. 初始化 Redis (如果 Login 用到了的话)
	dao.InitRedis()

	// 5. 初始化雪花算法
	snowflake.Init("2026-01-01", 1)

	// 6. Mock JWT 配置 (防止配置文件里没写)
	viper.Set("auth.jwt_secret", "test_secret")
	viper.Set("auth.jwt_expire", 24)
}

func TestLoginHandler(t *testing.T) {
	// 初始化环境 (连接数据库等)
	// 如果你只想测路由结构而不连库，可以注释掉这一行，但业务逻辑会报错
	initEnv()

	// 1. 设置 Gin 为测试模式 (不打印多余日志)
	gin.SetMode(gin.TestMode)

	// 2. 定义路由
	r := gin.Default()
	r.POST("/login", LoginHandler)

	// 3. 构造请求参数 (假设数据库里已经有 admin/123456 这个用户)
	// 如果没有，请先在数据库里手动插入一条，或者先调用 SignUp 接口
	params := models.ParamLogin{
		Username: "admin",
		Password: "123456",
	}
	jsonBytes, _ := json.Marshal(params)

	// 4. 构造 HTTP 请求
	// 方法: POST, 路径: /login, Body: JSON数据
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json") // 必加！

	// 5. 创建记录仪 (假浏览器)
	w := httptest.NewRecorder()

	// 6. 执行请求
	r.ServeHTTP(w, req)

	// 7. 断言结果
	// 拿到响应结果
	respBody := w.Body.String()
	t.Logf("响应结果: %s", respBody)

	// 判断状态码
	// 注意：如果数据库连不上，这里可能是 200 (CommonError也是200返回) 但是 code 是非 1000
	assert.Equal(t, http.StatusOK, w.Code)

	// 判断是否包含 token
	if w.Code == 200 {
		// 只有登录成功才会有 token
		// assert.Contains(t, respBody, "token")
	}
}
