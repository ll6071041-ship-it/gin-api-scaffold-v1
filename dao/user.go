package dao

import (
	"errors"
	"gin-api-scaffold-v1/models"

	"gorm.io/gorm"
)

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	// GORM 写法：Model指定查哪个表，Where指定条件，Count统计数量
	err = DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err // 数据库查询出错
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return nil
}

// InsertUser 插入新用户
func InsertUser(user *models.User) (err error) {
	// GORM 写法：直接 Create 一个结构体指针
	err = DB.Create(user).Error
	return
}

// GetUserByUsername 根据用户名查用户 (用于登录)
func GetUserByUsername(username string) (user *models.User, err error) {
	user = new(models.User)
	// GORM 写法：Where条件 + First(查第一条)
	err = DB.Where("username = ?", username).First(user).Error

	// 如果没查到，GORM 会返回 gorm.ErrRecordNotFound
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("用户不存在")
	}
	return
}
