package user

import (
	"github.com/gin-gonic/gin"
	common "mirey7/project-common"
)

type HandlerUser struct {
}

func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	rsp := &common.Result{}
	ctx.JSON(200, rsp.Success("123456"))
}
