package logic

import (
	"errors"
	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/models"
	"gin-api-scaffold-v1/pkg/encrypt"
	"gin-api-scaffold-v1/pkg/jwt" // 👈 1. 引入这一行
	"gin-api-scaffold-v1/pkg/snowflake"
)

// SignUp 处理注册业务 (保持不变)
func SignUp(p *models.ParamSignUp) (err error) {
	if err = dao.CheckUserExist(p.Username); err != nil {
		return err
	}
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: encrypt.EncryptPassword(p.Password),
	}
	return dao.InsertUser(user)
}

// Login 处理登录业务
func Login(p *models.ParamLogin) (token string, err error) {
	// 1. 去数据库查用户是否存在
	user, err := dao.GetUserByUsername(p.Username)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	// 2. 校验密码
	password := encrypt.EncryptPassword(p.Password)
	if password != user.Password {
		return "", errors.New("密码错误")
	}

	// 3. ⚡️⚡️ 生成标准的 JWT Token ⚡️⚡️
	// 使用我们在 pkg/jwt 里封装好的 GenToken 函数
	// 只要这一步不报错，前端拿到的就是一张合法的“通行证”
	return jwt.GenToken(user.UserID, user.Username)
}
