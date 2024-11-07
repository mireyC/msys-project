package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"mirey7/project-common/discovery"
	"mirey7/project-grpc/project"
	"mirey7/project-project/config"
	projectServiceV1 "mirey7/project-project/pkg/service/project.service.v1"
	"net"
)

// Router 接口
type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)

	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr        string
	RegisterFun func(*grpc.Server)
}

func ServerRegisterAndRun() {
	etcdV3Client := discovery.New(config.C.EC.Addr,
		config.C.EC.UserName,
		config.C.EC.Password,
		config.C.EC.DialTime)
	rgMtarget := discovery.RegisterManagerTarget(config.C.GC.Name)
	// windwos
	rgKey := discovery.RegisterKey(config.C.GC.Name, config.C.GC.Addr)
	rgVal := config.C.GC.Addr
	// linux
	//linuxRgKey := config.GetPublicIP() + config.C.GC.Port
	//rgKey := discovery.RegisterKey(config.C.GC.Name, linuxRgKey)
	//rgVal := linuxRgKey

	rgWeitht := config.C.GC.Weight
	err := etcdV3Client.Register(rgMtarget, rgKey, rgVal, rgWeitht)
	if err != nil {
		log.Fatalln("register etcd fail, cause by: ", err)
	}

	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFun: func(g *grpc.Server) {
			project.RegisterProjectServiceServer(g, projectServiceV1.New())
		},
	}
	s := grpc.NewServer()
	c.RegisterFun(s)

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("grpc server cannot listen cause by ", err)
	}

	//go func() {
	str := fmt.Sprintf("grpc server %s running on %s", config.C.GC.Name, config.C.GC.Addr)
	log.Println(str)
	err = s.Serve(lis)
	if err != nil {
		str := fmt.Sprintf("%s run on %s fail, cause by %v", config.C.GC.Name, config.C.GC.Addr, err)
		log.Fatalln(str)
		//log.Println("server started error ", er)
		//return
	}
	//}()

	//return s
}
