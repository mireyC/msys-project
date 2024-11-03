package model

import (
	"mirey7/project-common/errs"
)

var (
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)
