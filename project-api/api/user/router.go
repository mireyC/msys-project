package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"mirey7/project-api/router"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	ru := &RouterUser{}
	router.Register(ru)
}
func (*RouterUser) Route(r *gin.Engine) {
	// 初始化 grpc 的客户端连接
	InitRpcUserClient()
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.register)
	r.POST("project/login", h.login)
}
