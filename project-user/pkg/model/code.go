package model

import (
	"mirey7/project-common/errs"
)

var (
	RedisError          = errs.NewError(999, "Redis 错误")
	DBError             = errs.NewError(998, "DB 错误")
	NoLegalMobile       = errs.NewError(102001, "手机号不合法")
	CaptchaError        = errs.NewError(102002, "验证码错误")
	EmailExist          = errs.NewError(102003, "邮箱已存在")
	MobileExist         = errs.NewError(102004, "手机号已存在")
	AccountExist        = errs.NewError(102005, "账号已存在")
	CaptchaNoExist      = errs.NewError(102006, "验证码不存在或已过期 ")
	AccountAndPwdError  = errs.NewError(102007, "账号密码不正确")
	OrganizationNoFound = errs.NewError(102008, "组织查询不到")
	NoLogin             = errs.NewError(102009, "未登录")
)
