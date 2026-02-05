package jwt

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestJWT 是一个综合测试，涵盖了生成和解析
func TestJWT(t *testing.T) {
	// =========================================================================
	// 1. Mock 配置 (关键步骤！)
	// =========================================================================
	// 因为我们的 GenToken 和 ParseToken 依赖 viper 读取配置
	// 但单元测试不会去读 config.yaml，所以我们要手动设置假的配置项
	viper.Set("auth.jwt_secret", "my_test_secret_key") // 设置假的密钥
	viper.Set("auth.jwt_expire", 1)                    // 设置过期时间为 1 小时

	// =========================================================================
	// 2. 定义测试数据
	// =========================================================================
	userID := int64(10086)
	username := "qimi_test"

	// =========================================================================
	// 3. 测试生成 Token (GenToken)
	// =========================================================================
	token, err := GenToken(userID, username)

	// 断言：生成过程不应该报错
	if err != nil {
		t.Fatalf("GenToken failed: %v", err)
	}
	// 断言：生成的 Token 不应该是空字符串
	assert.NotEmpty(t, token, "Token 不应该为空")

	// =========================================================================
	// 4. 测试解析 Token (ParseToken)
	// =========================================================================
	claims, err := ParseToken(token)

	// 断言：解析过程不应该报错
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	// =========================================================================
	// 5. 验证数据一致性
	// =========================================================================
	// 核心：我们要验证解密出来的数据，是不是就是我们当初加密进去的
	assert.Equal(t, userID, claims.UserID, "UserID 应该一致")
	assert.Equal(t, username, claims.Username, "Username 应该一致")

	// 验证一下发行人
	assert.Equal(t, "gin-api-scaffold", claims.Issuer)

	t.Logf("JWT 测试通过！生成 Token: %s", token)
}

// TestParseInvalidToken 测试解析错误的 Token
func TestParseInvalidToken(t *testing.T) {
	viper.Set("auth.jwt_secret", "my_test_secret_key")

	// 随便瞎写一个 Token
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"

	claims, err := ParseToken(invalidToken)

	// 断言：这里必须报错
	assert.Error(t, err, "解析错误的 Token 应该返回错误")
	assert.Nil(t, claims, "解析失败 Claims 应该是 nil")
}
