package main

import "mirey7/project-user/router"

func main() {
	//r := gin.Default()

	// 路由
	//router.InitRouter(r)
	//gc := router.RegisterGrpc()
	//stop := func() {
	//	gc.Stop()
	//}

	//srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
	router.ServerRegisterAndRun()
}
