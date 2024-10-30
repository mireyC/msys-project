package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, serverName string, addr string) {
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 保证下面优雅启停
	go func() {
		log.Printf("%s running on %s \n", serverName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	// SIGINT 用户发送 INTR 字符（Ctrl+C）触发
	// SIGTERM 结束程序（可以被捕获）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting Down project %s...\n", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by: %v \n", serverName, err)
	}

	select {
	case <-ctx.Done():
		log.Println("关闭超时")
	}

	log.Printf("%s stop success...\n", serverName)
}
