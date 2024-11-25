package client

import (
	"mirey7/project-common/discovery"
	"mirey7/project-grpc/user/login"
	"mirey7/project-project/config"
)

var UserSvcClient login.LoginServiceClient

func InitUserSvcClient() {
	etcdV3Client := discovery.New(config.C.EC.Addr, config.C.EC.UserName, config.C.EC.Password, config.C.EC.DialTime)

	conn := etcdV3Client.InitRoundRobinGrpcConn(config.C.GC.UserService)
	UserSvcClient = login.NewLoginServiceClient(conn)
}
