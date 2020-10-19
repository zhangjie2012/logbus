package logbus

import (
	"context"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack"
	"github.com/zhangjie2012/cbl-go/cache"
)

// RedisListInput data from redis LIST
type RedisListInput struct {
	listKey string
}

func NewRedisListInput(app string, addr string, password string, db int, key string) (Input, error) {
	if err := cache.InitCache(app, addr, password, db); err != nil {
		return nil, err
	}
	return &RedisListInput{listKey: key}, nil
}

func (in *RedisListInput) Read(ctx context.Context) (*StandardLog, error) {

Loop:
	select {
	case <-ctx.Done():
		return nil, NoInputData
	default:
		bs, err := cache.MQBlockPop(in.listKey, 1*time.Second)
		if err != nil {
			if err == cache.NotExist {
				goto Loop
			}
			return nil, fmt.Errorf("read data error: %s", err)
		}

		l := StandardLog{}
		if err := msgpack.Unmarshal(bs, &l); err != nil {
			return nil, fmt.Errorf("marshal log error: %s", err)
		}
		return &l, nil
	}
}

func (in *RedisListInput) Close() error {
	cache.CloseCache()
	return nil
}
