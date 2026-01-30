package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/logger"
	"gin-api-scaffold-v1/routers"
	"gin-api-scaffold-v1/settings"
)

func main() {
	// 1. 加载配置
	if err := settings.InitConfig(); err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 2. 初始化日志
	logger.InitLogger()
	defer logger.Logger.Sync()

	// 3. 初始化 MySQL
	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}

	// 4. 初始化 Redis
	if err := dao.InitRedis(); err != nil {
		panic(err)
	}

	// 5. 注册路由
	r := routers.SetupRouter()

	// 6. 启动服务 (优雅关机版)
	// ==========================================
	// 我们不再使用 r.Run()，因为那个太简陋了
	// 我们要自己定义一个 http.Server
	// ==========================================

	port := viper.GetString("app.port")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 【开启一个协程启动服务】
	// 为什么要用协程？因为 srv.ListenAndServe() 会卡死在这里不动
	// 我们需要它在后台运行，主线程要空出来去监听“关机信号”
	go func() {
		fmt.Printf("服务正在启动，端口: %s\n", port)
		// 这里的 ErrServerClosed 是正常的关闭信号，不是错误
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("listen: ", zap.Error(err))
		}
	}()

	// 【等待关机信号】
	// 创建一个通道，专门等系统发信号
	quit := make(chan os.Signal, 1)

	// 告诉系统：如果有 SIGINT (Ctrl+C) 或 SIGTERM (Docker杀进程)，请通知我
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 程序会卡在这里死等，直到收到信号
	<-quit
	logger.Logger.Info("Shutdown Server ...")

	// 【优雅关机】
	// 创建一个 5 秒的超时上下文
	// 意思是：我给服务器 5 秒钟时间把手头的事情做完，5秒后不管做没做完都强制关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Shutdown 通知服务器停止接客，处理剩下的请求
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Server Shutdown:", zap.Error(err))
	}

	logger.Logger.Info("Server exiting")
}
