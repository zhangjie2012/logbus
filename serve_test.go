package logbus

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack"
	"github.com/zhangjie2012/cbl-go/cache"
)

func TestServe(t *testing.T) {
	logKey := "TestRedisListInput.Log"
	appName := "TestRedisListInput"
	timestamp := time.Now().UnixNano()
	level := logrus.DebugLevel
	msg := "hello"

	// fake data
	go func() {
		l := StdLog{
			AppName:     appName,
			Timestamp:   timestamp,
			Level:       level.String(),
			StateId:     StateIdInvalid,
			Caller:      "",
			Message:     msg,
			Annotations: map[string]interface{}{"hello": "world"},
		}
		bs, err := msgpack.Marshal(&l)
		require.Nil(t, err)

		for i := 0; i < 10; i++ {
			cache.MQPush(logKey, bs)
			time.Sleep(1 * time.Second)
		}
	}()

	in, err := NewRedisListInput(appName, "localhost:6379", "", 0, logKey)
	require.Nil(t, err)

	out, err := NewStdoutOutput(DefaultTransformer)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go Serve(ctx, in, []Output{out})

	time.Sleep(10 * time.Second)
}
