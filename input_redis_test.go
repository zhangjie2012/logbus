package logbus

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack"
	"github.com/zhangjie2012/cbl-go/cache"
)

func TestRedisListInput(t *testing.T) {
	logKey := "TestRedisListInput.Log"
	appName := "TestRedisListInput"
	timestamp := time.Now().UnixNano()
	level := logrus.DebugLevel
	msg := "hello"

	input, err := NewRedisListInput(appName, "localhost:6379", "", 0, logKey)
	require.Nil(t, err)

	go func() {
		l := StandardLog{
			AppName:   appName,
			Timestamp: timestamp,
			Level:     level.String(),
			StateId:   StateIdInvalid,
			Caller:    "",
			Message:   msg,
		}
		bs, err := msgpack.Marshal(&l)
		require.Nil(t, err)

		for i := 0; i < 10; i++ {
			cache.MQPush(logKey, bs)
			time.Sleep(3 * time.Second)
		}
	}()

	count := 0
	for count < 10 {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			l, err := input.Read(ctx)
			if err == NoInputData {
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, msg, l.Message)
			count++
		}()
	}
}
