package dao

import (
	"fmt"

	"github.com/spf13/viper" // 👈 引入 viper
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err // ✅ 直接返回连接结果，不要去建表
}
