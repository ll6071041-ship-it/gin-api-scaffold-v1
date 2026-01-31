package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器，方便在其他地方调用
var Trans ut.Translator

// InitTrans 初始化翻译器
// locale: 语言环境，通常传 "zh"
func InitTrans(locale string) (err error) {
	// 1. 修改 Gin 框架中的 Validator 引擎属性，实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// =============================================================
		// 🔥 核心功能：自定义错误字段名 (使用 json tag)
		// =============================================================
		// 注册一个获取 json tag 的自定义方法
		// 默认情况下 validator 返回的是结构体字段名 (如 "UserName")
		// 这样写之后，校验失败时，错误信息就会显示 "user_name" (即 json tag 的值)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 获取 tag 中的 json 值，例如 `json:"user_name,omitempty"`
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// 如果 json tag 是 "-"，说明忽略该字段，返回空
			if name == "-" {
				return ""
			}
			return name
		})

		// 2. 初始化翻译器
		zhT := zh.New() // 中文翻译器
		// 第一个参数是备用语言，第二个是当前语言
		uni := ut.New(zhT, zhT)

		// 获取具体的翻译实例
		// 通常我们 locale 传 "zh"，这里就会获取到中文翻译器
		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// =============================================================
		// 🔥 核心功能：注册中文翻译
		// =============================================================
		// 这一步把 validator 内置的英文错误信息替换成中文
		switch locale {
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(v, Trans)
		default:
			err = zh_translations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

// RemoveTopStruct 去除结构体名称前缀
// validator 返回的错误 key 默认是 "StructName.FieldName" (例如 "SignUpParam.Password")
// 我们想要的是纯粹的 "password" 或者 "mobile"
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		// field 可能是 "SignUpParam.password"
		// err 是翻译后的错误信息，例如 "password 为必填字段"

		// 截取点号之后的部分
		// strings.Index(field, ".") 返回点号的位置
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
