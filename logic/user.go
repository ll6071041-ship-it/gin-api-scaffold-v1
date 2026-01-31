package logic

import (
	// ❌ 删掉原来的 "gin-api-scaffold-v1/dao/mysql"
	// ✅ 改成下面这个，因为你的 dao 层都在 package dao 下
	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/models"
	"gin-api-scaffold-v1/pkg/encrypt"
	"gin-api-scaffold-v1/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户存不存在 (调用刚才在 dao/user.go 里写的函数)
	if err = dao.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成 UserID
	userID := snowflake.GenID()

	// 3. 构造 User 实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: encrypt.EncryptPassword(p.Password), // 密码加密
	}

	// 4. 保存进数据库
	return dao.InsertUser(user)
}
