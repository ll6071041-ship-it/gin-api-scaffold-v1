package dao

import (
	"errors"
	"gin-api-scaffold-v1/models"

	"gorm.io/gorm"
)

// =================================================================
// ⚡️ 1. 定义全局错误变量 (给 Controller 层判断逻辑用的)
// =================================================================
var (
	ErrorUserExist    = errors.New("用户已存在")
	ErrorUserNotFound = errors.New("用户不存在")
)

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	err = DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		// ⚡️ 2. 这里返回上面定义的全局变量，而不是临时 new 一个
		return ErrorUserExist
	}
	return nil
}

// InsertUser 插入新用户
func InsertUser(user *models.User) (err error) {
	err = DB.Create(user).Error
	return
}

// GetUserByUsername 根据用户名查用户 (用于登录)
func GetUserByUsername(username string) (user *models.User, err error) {
	user = new(models.User)
	err = DB.Where("username = ?", username).First(user).Error

	if err == gorm.ErrRecordNotFound {
		// ⚡️ 3. 这里也返回全局变量
		return nil, ErrorUserNotFound
	}
	return
}
