package repo

/*
	repo
	封装数据库操作的逻辑，提供应用层访问数据库的接口
*/

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, val string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
