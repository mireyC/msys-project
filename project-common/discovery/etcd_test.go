package discovery

import (
	"context"
	"fmt"
	etcdV3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := etcdV3.New(etcdV3.Config{
		Endpoints:   []string{"120.26.242.2:2379"},
		Username:    "root",          // 使用之前设置的用户名
		Password:    "mirey7/A",      // 使用之前设置的密码
		DialTimeout: 5 * time.Second, // 设置连接超时
	})
	if err != nil {
		log.Fatal("etcd client init fail, cause by: ", err)
	}

	// 使用Put操作
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 将键值对 "my-key" -> "my-value" 存入etcd
	_, err = client.Put(ctx, "ccs", "my-value")
	if err != nil {
		log.Fatalf("Failed to put key-value pair: %v", err)
	}

	// 检查存储结果
	resp, err := client.Get(ctx, "cc")
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}

	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", kv.Key, kv.Value)
	}
}
