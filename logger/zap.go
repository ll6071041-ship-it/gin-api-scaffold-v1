package logger

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义一个全局变量，方便内部调用（虽然有了 zap.L()，但保留这个是个好习惯）
var Logger *zap.Logger

// InitLogger 初始化 Logger (企业级完整版)
// 负责把日志系统跑起来，配置好“写到哪里”、“怎么写”、“记哪些级别”
func InitLogger() {
	// =================================================================
	// 1. 获取日志写入器 (Writer)
	// =================================================================
	// 我们不直接用 os.OpenFile，因为文件会越来越大。
	// 这里用 lumberjack 库，它能自动“切割”日志文件（比如达到 10MB 就换个新文件）。
	writeSyncer := getLogWriter(
		viper.GetString("log.filename"), // 从配置读取文件名 (如: ./logs/bluebell.log)
		viper.GetInt("log.max_size"),    // 单个文件最大尺寸 (MB)
		viper.GetInt("log.max_backups"), // 最多保留几个旧文件
		viper.GetInt("log.max_age"),     // 旧文件最多保留几天
	)

	// =================================================================
	// 2. 获取日志编码器 (Encoder)
	// =================================================================
	// 决定日志长什么样。是 {"msg":"hello"} 这种 JSON 格式？
	// 还是 [INFO] 2023-01-01 hello 这种普通文本格式？
	encoder := getEncoder()

	// =================================================================
	// 3. 定义日志级别 (Level)
	// =================================================================
	// 从配置文件读 log.level (比如 "debug", "info", "error")
	// 只有大于等于这个级别的日志才会被记录。
	var l = new(zapcore.Level)
	if err := l.UnmarshalText([]byte(viper.GetString("log.level"))); err != nil {
		// 如果配置文件填错了，默认给个 Debug 级别，保证能打出日志
		*l = zapcore.DebugLevel
	}

	// =================================================================
	// 4. 创建 Core (核心引擎)
	// =================================================================
	// Core 是 Zap 的心脏，它把上面三者结合在一起：
	// Encoder: 怎么编码
	// WriteSyncer: 写到哪里
	// Level: 记什么级别

	// ⚡️ NewMultiWriteSyncer: 这是一个神器。
	// 它能让日志“分身”，同时写到两个地方：
	// 1. writeSyncer -> 写到日志文件里 (持久化保存)
	// 2. os.Stdout   -> 写到黑窗口/控制台里 (方便开发时实时看)
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)),
		l,
	)

	// =================================================================
	// 5. 构造 Logger 对象
	// =================================================================
	// zap.AddCaller(): 非常重要！
	// 加上它，日志里就会显示是哪个文件、哪一行打印的 (例如: main.go:15)，
	// 否则你出了 Bug 根本找不到在哪里。
	Logger = zap.New(core, zap.AddCaller())

	// =================================================================
	// 6. ⚡️⚡️⚡️ 替换全局 Logger ⚡️⚡️⚡️
	// =================================================================
	// 这一步是关键！
	// Zap 默认有一个全局 Logger，但那个是空的，不干活。
	// 我们把自己辛苦配置好的 Logger 替换上去。
	// 以后你在任何地方调用 zap.L().Info()，实际上用的就是我们配置好的这个 Logger。
	zap.ReplaceGlobals(Logger)
}

// ---------------------------------------------------------------------
// 内部辅助函数
// ---------------------------------------------------------------------

// getLogWriter 配置日志切割规则
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件的位置
		MaxSize:    maxSize,    // 文件到多大就开始切割 (MB)
		MaxBackups: maxBackups, // 切割后，保留几个旧文件 (防止硬盘被日志占满)
		MaxAge:     maxAge,     // 旧文件保留多少天
		Compress:   false,      // 是否压缩旧文件 (gzip)，一般设为 false 方便直接看
	}
	// AddSync 把 lumberjack 转成 zap 需要的 WriteSyncer 类型
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoder 配置日志的格式
func getEncoder() zapcore.Encoder {
	// 使用生产环境的默认配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// ⚡️ 修改时间格式：
	// 默认是时间戳 (1678888.123)，人类看不懂。
	// 改成 ISO8601 (2023-01-01T12:00:00.000Z)，一眼就能看懂时间。
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// ⚡️ 修改级别显示：
	// 把 info 变成 INFO，debug 变成 DEBUG，醒目一点。
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 返回一个 ConsoleEncoder (控制台格式/普通文本格式)
	// 这种格式带颜色，人看着舒服。
	// 如果是上线跑在 Docker/K8s 里收集日志，通常建议换成 zap.NewJSONEncoder(encoderConfig)
	return zapcore.NewConsoleEncoder(encoderConfig)
}
