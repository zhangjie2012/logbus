package logbus

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack"
)

// RedisListInput data from redis LIST
type RedisListInput struct {
	client  *redis.Client
	listKey string
}

func NewRedisListInput(app string, addr string, password string, db int, key string) (Input, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return &RedisListInput{client: client, listKey: key}, nil
}

func (in *RedisListInput) Read(ctx context.Context) (*StdLog, error) {

Loop:
	select {
	case <-ctx.Done():
		return nil, NoInputData
	default:
		result, err := in.client.BLPop(1*time.Second, in.listKey).Result()
		if err != nil && err != redis.Nil {
			return nil, fmt.Errorf("read data error: %s", err)
		}
		if err == redis.Nil {
			goto Loop
		}

		// result[0] is key name
		bs := []byte(result[1])
		l := StdLog{}
		if err := msgpack.Unmarshal(bs, &l); err != nil {
			return nil, fmt.Errorf("marshal log error: %s", err)
		}
		return &l, nil
	}
}

func (in *RedisListInput) Close() error {
	return in.client.Close()
}
