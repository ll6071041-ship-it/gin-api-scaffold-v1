package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// secret 是一个盐值，随便写什么都行，保密级别越高越好
// 它的作用是：假设用户的密码是 "123456"，加上盐变成 "123456woshimima"，
// 这样生成的 MD5 值就和普通的 "123456" 不一样了，黑客猜不到。
const secret = "bluebell.app.secret.key"

// EncryptPassword 将明文密码加密为 MD5 字符串
func EncryptPassword(oPassword string) string {
	h := md5.New()
	// 把 盐 + 密码 拼接起来一起加密
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
