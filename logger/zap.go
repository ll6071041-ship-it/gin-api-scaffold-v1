package logger // 👈 注意这里是 logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局变量，给其他地方用
var Logger *zap.Logger

// InitLogger 初始化 Logger
func InitLogger() {
	// ... 这里粘贴我之前给你的 "万能配置模板" 里的 InitLogger 函数体代码 ...
	// 为了省篇幅，核心逻辑就是：配置 Encoder -> 配置 WriteSyncer -> zap.New()

	// 简略版示例（你可以替换成之前详细版的）：
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	file, _ := os.Create("./my-app.log")
	writeSyncer := zapcore.AddSync(file)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	Logger = zap.New(core, zap.AddCaller())
}
