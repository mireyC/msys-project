package dao

/*
	dao
	提供应用层访问数据库接口的具体实现
*/

import (
	"context"
	"github.com/go-redis/redis/v8"
	"mirey7/project-user/config"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	rdb := redis.NewClient(config.C.ReadRedisConfig())

	Rc = &RedisCache{
		rdb: rdb,
	}
}

func (rc *RedisCache) Put(ctx context.Context, key, val string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, val, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.rdb.Get(ctx, key).Result()
	return res, err
}
