package logic

import (
	"errors"
	"gin-api-scaffold-v1/dao"
	"gin-api-scaffold-v1/models"
	"gin-api-scaffold-v1/pkg/encrypt"
	"gin-api-scaffold-v1/pkg/snowflake"
)

// SignUp 处理注册业务
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户存不存在
	if err = dao.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成 UserID
	userID := snowflake.GenID()

	// 3. 构造 User 实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: encrypt.EncryptPassword(p.Password), // ⚡️ 密码加密存储
	}

	// 4. 保存进数据库
	return dao.InsertUser(user)
}

// Login 处理登录业务
func Login(p *models.ParamLogin) (token string, err error) {
	// 1. 去数据库查用户是否存在
	user, err := dao.GetUserByUsername(p.Username)
	if err != nil {
		// 查不到用户，直接返回错误
		return "", errors.New("用户不存在")
	}

	// 2. ⚡️ 校验密码
	// 逻辑：把前端传来的密码(p.Password)进行同样的加密，然后跟数据库里的密文(user.Password)比对
	password := encrypt.EncryptPassword(p.Password)
	if password != user.Password {
		return "", errors.New("密码错误")
	}

	// 3. 生成 Token
	// 目前我们还没集成 JWT，先返回一个假的 Token 占位，保证代码能跑通
	// 下一步我们可以引入 JWT
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fake.token", nil
}
