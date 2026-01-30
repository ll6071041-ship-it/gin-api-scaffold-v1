package logger

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化 Logger (企业级完整版)
func InitLogger() {
	// 1. 获取日志写入器 (用 lumberjack 实现切割)
	writeSyncer := getLogWriter(
		viper.GetString("log.filename"), // 从配置读取文件名
		viper.GetInt("log.max_size"),    // max_size MB
		viper.GetInt("log.max_backups"), // 保留几个旧文件
		viper.GetInt("log.max_age"),     // 保留几天
	)

	// 2. 获取日志编码器 (让日志变成 JSON 还是 普通文本)
	encoder := getEncoder()

	// 3. 定义日志级别
	// 在生产环境可以是 InfoLevel，开发环境可以是 DebugLevel
	var l = new(zapcore.Level)
	if err := l.UnmarshalText([]byte(viper.GetString("log.level"))); err != nil {
		// 默认级别
		*l = zapcore.DebugLevel
	}

	// 4. 创建 Core
	// NewCore 接收三个参数：编码器、写入器、日志级别
	// ⚡️这是重点：我们要同时输出到 文件 和 控制台
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), // 同时写文件和控制台
		l,
	)

	// 5. 构造 Logger
	// AddCaller: 显示文件名和行号
	Logger = zap.New(core, zap.AddCaller())
}

// 辅助函数 1: 配置日志切割
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件最大尺寸 (MB)
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		MaxAge:     maxAge,     // 保留旧文件的最大天数
		Compress:   false,      // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 辅助函数 2: 配置编码器 (让日志好看一点)
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	// ⚡️修改时间格式：把 167888.123 变成 2023-01-01 12:00:00
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// ⚡️日志级别大写：info -> INFO
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 如果你喜欢 JSON 格式 (机器读)，用 NewJSONEncoder
	// 如果你喜欢 普通文本格式 (人读)，用 NewConsoleEncoder
	// 这里建议：开发阶段用 Console，上线用 JSON。为了方便你现在看，我先用 Console
	return zapcore.NewConsoleEncoder(encoderConfig)
}
