package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"mirey7/project-api/api/midd"
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
	group := r.Group("/project/index")
	group.Use(midd.TokenVerify())
	group.POST("", h.Index)
	group1 := r.Group("/project/project")
	group1.Use(midd.TokenVerify())
	group1.POST("/selfList", h.myProjectList)
	group1.POST("", h.myProjectList)
}
