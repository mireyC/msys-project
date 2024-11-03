package main

import (
	"github.com/gin-gonic/gin"
	_ "mirey7/project-api/api"
	"mirey7/project-api/config"
	"mirey7/project-api/router"
	srv "mirey7/project-common"
)

func main() {
	r := gin.Default()

	// 路由
	router.InitRouter(r)

	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}
