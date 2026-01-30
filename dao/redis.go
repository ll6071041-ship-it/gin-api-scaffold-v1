package dao

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RDB 作为一个全局的客户端连接实例，给同一个包下的 redis_todo.go 使用
var RDB *redis.Client

// InitRedis 初始化连接
func InitRedis() (err error) {
	RDB = redis.NewClient(&redis.Options{
		// 拼接地址：IP:Port
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	// 测试一下连接
	_, err = RDB.Ping(context.Background()).Result()
	return err
}
