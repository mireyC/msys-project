package project

import (
	"mirey7/project-api/config"
	"mirey7/project-common/discovery"
	"mirey7/project-grpc/project"
)

var ProjectSvcClient project.ProjectServiceClient

func InitRpcProjectClient() {

	//log.Printf("grpc client: %s", target)
	//conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	etcdV3Client := discovery.New(config.C.EC.Addr, config.C.EC.UserName, config.C.EC.Password, config.C.EC.DialTime)

	conn := etcdV3Client.InitRoundRobinGrpcConn(config.C.GC.ProjectService)
	ProjectSvcClient = project.NewProjectServiceClient(conn)
}
