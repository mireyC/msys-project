package main

import (
	"github.com/gin-gonic/gin"
	srv "mirey7/project-common"
	_ "mirey7/project-user/api"
	"mirey7/project-user/router"
)

func main() {
	r := gin.Default()
	// 路由
	router.InitRouter(r)

	srv.Run(r, "project-user", ":80")
}
