package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitConfig 读取配置文件
func InitConfig() error {
	viper.SetConfigName("config") // 文件名 (不带后缀)
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 搜索路径 (当前目录)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置失败: %s", err)
	}

	// 开启热加载 (改了配置文件不用重启)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
	})

	return nil
}
