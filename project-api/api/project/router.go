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
	group := r.Group("/project")
	group.Use(midd.TokenVerify())
	group.POST("/index", h.Index)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.projectRead)
	group.POST("/project/recycle", h.projectRecycle)
	group.POST("/delProject", h.delProject)
	group.POST("/project_collect/collect", h.projectCollect)
	group.POST("/project/recovery", h.projectRecovery)
	group.POST("/project/edit", h.projectEdit)
}
