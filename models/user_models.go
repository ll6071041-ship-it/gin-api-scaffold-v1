package models

// SignUpParam 用户注册请求参数
type SignUpParam struct {
	// 1. 基础校验：必填
	// json:"username" -> 报错时显示 "username"
	// binding:"required" -> 必填
	Username string `json:"username" binding:"required"`

	// 2. 长度与格式校验
	// min=6,max=20 -> 长度限制在 6 到 20 之间
	Password string `json:"password" binding:"required,min=6,max=20"`

	// 3. 跨字段校验 (最酷的功能！)
	// eqfield=Password -> 这个字段的值必须和 Password 字段的值完全一样
	// 常用于“确认密码”场景
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`

	// 4. 类型校验
	// email -> 必须是合法的邮箱格式 (xxx@xxx.com)
	Email string `json:"email" binding:"required,email"`

	// 5. 数值范围校验
	// gte=1 -> 大于等于 1 (Greater Than or Equal)
	// lte=130 -> 小于等于 130 (Less Than or Equal)
	Age uint8 `json:"age" binding:"gte=1,lte=130"`

	// 6. 枚举校验 (Enum)
	// oneof -> 值必须是列出的其中一个，只能是 male, female 或 secret
	Gender string `json:"gender" binding:"required,oneof=male female secret"`

	// 7. 可选字段校验 (omitempty)
	// omitempty -> 如果前端没传这个字段（或者传空值），就不校验
	// url -> 如果传了，就必须是合法的 URL 格式 (http://...)
	PersonalSite string `json:"personal_site" binding:"omitempty,url"`
}

// RegisterParam 用户注册请求参数
// 这里的结构体名 RegisterParam 必须首字母大写，才能被其他包（如 controller）访问
type RegisterParam struct {
	// binding:"required": 必填
	// json:"username":    报错时我们希望显示 "username" 而不是 "Username"
	Username string `json:"username" binding:"required"`

	// binding:"gte=6":    长度必须大于等于 6
	Password string `json:"password" binding:"required,gte=6"`

	// binding:"email":    必须是合法的邮箱格式
	Email string `json:"email" binding:"required,email"`

	// binding:"gte=18":   数值必须大于等于 18
	// binding:"lte=130":  数值必须小于等于 130
	Age uint8 `json:"age" binding:"gte=18,lte=130"`

	// 演示非必填字段
	Address string `json:"address"`
}
