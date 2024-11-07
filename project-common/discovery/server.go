package discovery

import (
	"context"
	"fmt"
	etcdV3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mirey7/project-common/logs"
	"time"
)

type EtcdClient struct {
	client *etcdV3.Client
}

func New(addr, userName, password string, dialTime int64) *EtcdClient {
	client, err := etcdV3.New(etcdV3.Config{
		Endpoints:   []string{addr},
		Username:    userName,
		Password:    password,
		DialTimeout: time.Duration(dialTime) * time.Second,
	})
	if err != nil {
		log.Fatal("etcd client init fail , cause by: ", err)
	}
	return &EtcdClient{
		client: client,
	}
}

// RegisterKey
// service/user/192.168.1.1:50051
func RegisterKey(serviceName, addr string) (key string) {
	key = fmt.Sprintf("service/%s/%s", serviceName, addr)
	log.Println("rg key: ", key)
	return key
}

// RegisterManagerTarget
// service/user
func RegisterManagerTarget(serviceName string) (target string) {
	target = fmt.Sprintf("service/%s", serviceName)
	fmt.Println("rg target: ", target)
	fmt.Println("serviceName: ", serviceName)
	return target
}

// Register 服务注册
// mangerTarget / eg. service/user
// key => grpc的 服务名 + addr / eg. service/user/192.168.1.1:50051
// val => grpc监听的 addr / eg. 192.168.1.1:50051
// weight => 节点权重
func (c *EtcdClient) Register(mangerTarget, key, value string, weight int64) error {
	em, err := endpoints.NewManager(c.client, mangerTarget)
	if err != nil {
		log.Fatal("newManager err cause by: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var ttl int64 = 60
	leaseResp, err := c.client.Grant(ctx, ttl)
	if err != nil {
		log.Fatal("etcd client.Grant err, ", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: value,
		Metadata: map[string]any{
			"weight": weight,
			//"cpu":    90,
		},
	}, etcdV3.WithLease(leaseResp.ID))

	if err != nil {
		log.Fatal("etcd AddEndpoint err, ", err)
	}

	go func() {
		kaCtx, _ := context.WithCancel(context.Background())
		_, er := c.client.KeepAlive(kaCtx, leaseResp.ID)
		if er != nil {
			msg := fmt.Sprintf("keepAlive fail addr:%s", value)
			logs.LG.Error(msg)
		}
	}()

	return err
}

func loggingInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	fmt.Printf("Calling %s on target: %s\n", method, cc.Target())
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (c *EtcdClient) InitRoundRobinGrpcConn(serviceName string) *grpc.ClientConn {

	// 指定服务路径（按你服务注册的路径来，比如 "service/user"）
	servicePath := "service/user"

	// 获取该路径下的所有键值对
	resp, err := c.client.Get(context.Background(), servicePath, etcdV3.WithPrefix())
	if err != nil {
		log.Fatal("Failed to get service entries from etcd:", err)
	}
	// 输出所有服务实例的 IP 地址
	fmt.Println("Registered IP addresses for service:", servicePath)
	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", kv.Key, kv.Value)
	}

	bd, err := resolver.NewBuilder(c.client)
	if err != nil {
		log.Fatalln("resolver NewBuilder err: ", err)
	}

	svcCfg := `
{
    "loadBalancingConfig": [
        {
            "round_robin": {}
        }
    ]
}
`
	target := fmt.Sprintf("etcd:///service/%s", serviceName)
	// 修改 grpc.Dial 的调用，添加拦截器
	log.Printf("grpc client dial on %s \n", target)
	cc, err := grpc.Dial(target,
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(svcCfg),         // 添加负载均衡
		grpc.WithUnaryInterceptor(loggingInterceptor), // 添加拦截器
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("grpc client conn err, ", err)
	}

	return cc
}

// Resolve 服务发现
func Resolve() {

}
