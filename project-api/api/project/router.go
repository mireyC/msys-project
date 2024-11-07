package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"mirey7/project-api/router"
)

type RouterProject struct {
}

func init() {
	log.Println("init project router")
	ru := &RouterProject{}
	router.Register(ru)
}
func (*RouterProject) Route(r *gin.Engine) {
	// 初始化 grpc 的客户端连接
	InitRpcProjectClient()
	h := New()
	r.POST("/project/index", h.Index)

}
