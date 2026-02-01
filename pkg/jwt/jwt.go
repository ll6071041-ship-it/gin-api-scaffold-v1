package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// 对应图片里的 CustomClaims
// 我们需要把 UserID 和 Username 存到 Token 里面，这样后端就不需要查数据库也能知道是谁
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT
// 对应图片里的 GenToken 函数
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间：从配置文件读取 (比如 24 小时)
			// 相比图片里的硬编码 time.Hour * 24，这样更灵活
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour)),
			// 签发人
			Issuer: "gin-api-scaffold",
		},
	}
	// 使用指定的签名方法创建签名对象 (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的 Secret 签名并获得完整的编码后的字符串 Token
	// ⚡️ 注意：这里千万不要像图片里那样写死 []byte("夏天夏天...")
	// 而是从配置文件读取，保证安全
	return token.SignedString([]byte(viper.GetString("auth.jwt_secret")))
}

// ParseToken 解析 JWT
// 对应图片里的 ParseToken 函数
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// 返回签名用的密钥
		return []byte(viper.GetString("auth.jwt_secret")), nil
	})
	if err != nil {
		return nil, err
	}
	// 校验 token 是否有效
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
