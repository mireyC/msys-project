package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"mirey7/project-user/router"
)

func init() {
	log.Println("init user router...")
	router.Register(&RouterUser{})
}

type RouterUser struct {
}

func (*RouterUser) Route(r *gin.Engine) {
	h := &HandlerUser{}
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
