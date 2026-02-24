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
	// ⚠️ 注意：这里必须引入 docs 包，否则 Swagger 无法加载文档数据
	_ "gin-api-scaffold-v1/docs"
	"gin-api-scaffold-v1/logger"
	"gin-api-scaffold-v1/pkg/rabbitmq"
	"gin-api-scaffold-v1/pkg/snowflake"

	// 👇 引入我们刚刚写的 validator 包，起个别名 myValidator 防止和官方包重名
	myValidator "gin-api-scaffold-v1/pkg/validator"
	"gin-api-scaffold-v1/routers"
	"gin-api-scaffold-v1/settings"
)

// @title           Bluebell项目接口文档
// @version         1.0
// @description     这是一个基于Gin框架的社区后端项目(仿Reddit)
// @termsOfService  http://swagger.io/terms/

// @contact.name    KKK
// @contact.url     http://www.liwenzhou.com
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// main 函数是 Go Web 项目的唯一入口
func main() {
	// =========================================================================
	// 1. 加载配置 (Viper)
	// =========================================================================
	// 这一步读取 config.yaml 文件。如果连配置文件都读不到（比如文件不存在或格式错误），
	// 后续代码无法运行，所以直接 Panic 终止程序是合理的。
	if err := settings.InitConfig(); err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// =========================================================================
	// 2. 初始化日志 (Zap)
	// =========================================================================
	// ⚠️ 这一步必须尽早执行！
	// 只有执行了 InitLogger，全局的 zap.L() 才会被配置好。
	// 如果不初始化，后续调用 zap.L().Info() 将不会输出任何内容。
	logger.InitLogger()

	// defer Sync(): 这是一个好习惯。
	// 它的作用是：在 main 函数结束前，把内存里缓存的日志强制刷入硬盘。
	// 防止程序崩溃或退出时，最后几条关键日志丢失。
	defer zap.L().Sync()

	// 打印一条 Debug 日志，确认日志系统启动成功
	zap.L().Debug("logger init success...")

	// =========================================================================
	// 3. 初始化雪花算法 (Snowflake)
	// =========================================================================
	// 用于生成全局唯一的 ID (比如订单号、用户ID)。
	// 参数1 "2026-01-01": 项目的起始时间。
	// 参数2 1: 当前机器 ID (MachineID)。
	// ⚠️ 注意：如果是分布式部署（多台服务器），每台机器的 MachineID 必须不同！
	if err := snowflake.Init("2026-01-01", 1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// =========================================================================
	// 🔥 4. 新增：初始化 Validator 翻译器 (汉化校验)
	// =========================================================================
	// 这一步加载中文语言包，并注册 json tag 自定义方法。
	// 如果失败，意味着参数校验返回的错误全是英文且格式混乱，严重影响前端体验，所以建议处理错误。
	if err := myValidator.InitTrans("zh"); err != nil {
		fmt.Printf("init validator failed, err:%v\n", err)
		return
	}

	// =========================================================================
	// 5. 初始化 MySQL (GORM)
	// =========================================================================
	// 建立数据库连接池。如果连不上数据库，后端服务没有意义，直接 Panic。
	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}
	// GORM 自身维护连接池，通常不需要像 sql.DB 那样手动 defer Close

	// =========================================================================
	// 6. 初始化 Redis
	// =========================================================================
	// 建立 Redis 连接池，用于缓存或 session 管理。
	if err := dao.InitRedis(); err != nil {
		panic(err)
	}

	rabbitmq.Init()
	defer rabbitmq.Conn.Close()
	defer rabbitmq.Channel.Close()

	// 2. 启动消费者（让它在后台一直监听）
	rabbitmq.StartConsumer()

	// =========================================================================
	// 7. 注册路由 (Gin)
	// =========================================================================
	// 加载 routers 包里定义的所有 URL 规则、中间件（日志、跨域、Recovery）。
	r := routers.SetupRouter()

	// =========================================================================
	// 8. 启动服务 (HTTP Server)
	// =========================================================================
	port := viper.GetString("app.port") // 从配置文件读取端口

	// 手动创建一个 http.Server，而不是直接用 r.Run()
	// 原因：r.Run() 内部也是创建 server，但它不方便我们在外部调用 Shutdown 做优雅关机。
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 🚀 开启一个协程 (Goroutine) 启动服务
	// 为什么？因为 srv.ListenAndServe() 是一个“死循环”，它会卡住当前线程一直等待请求。
	// 如果不放在协程里，主线程卡在这里，后面的“优雅关机”代码永远执行不到。
	go func() {
		fmt.Printf("服务正在启动，端口: %s\n", port)
		zap.L().Info("Server is starting...", zap.String("port", port))

		// 启动服务。
		// 如果返回错误，且错误不是“服务器已关闭(http.ErrServerClosed)”，说明启动失败（比如端口被占用）。
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Fatal 会打印错误日志并直接退出程序 (os.Exit)
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// =========================================================================
	// 9. 优雅关机 (Graceful Shutdown)
	// =========================================================================

	// 创建一个通道，专门用来接收操作系统的信号
	quit := make(chan os.Signal, 1)

	// signal.Notify 告诉操作系统：
	// 如果收到了 SIGINT (Ctrl+C) 或者 SIGTERM (Docker 停止容器/K8s 销毁 Pod)，
	// 请把信号发给 quit 通道，不要直接把我的程序杀掉！让我有时间处理后事。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 🛑 程序会卡在这里“死等”！
	// 直到 quit 通道里收到了信号，主线程才会继续往下执行。
	<-quit

	zap.L().Info("Shutdown Server ...")

	// 创建一个 5 秒的超时上下文
	// 意思是：我给服务器 5 秒钟的时间去处理手里还没处理完的请求。
	// 如果 5 秒到了还没处理完，就强制关闭，不再等了。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// srv.Shutdown 会做两件事：
	// 1. 马上停止接收新的请求。
	// 2. 等待正在处理的请求处理完（或者直到 ctx 超时）。
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
