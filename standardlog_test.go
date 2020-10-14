package logbus

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardLog(t *testing.T) {
	appName := "logbus"
	now := time.Now()
	metadata := logrus.Fields{"stateid": "register", "ip": "127.0.0.1", "userid": "zhangjie2012"}
	level := logrus.DebugLevel
	message := "hello, world"

	l := newStandardLog(appName, now, metadata, nil, level, message)
	require.NotNil(t, l)

	assert.Equal(t, appName, l.AppName)
	assert.EqualValues(t, now.UnixNano(), l.Timestamp)
	assert.EqualValues(t, level.String(), l.Level)
	assert.EqualValues(t, message, l.Message)
	assert.EqualValues(t, "register", l.StateId)
	assert.EqualValues(t, map[string]interface{}{"ip": "127.0.0.1", "userid": "zhangjie2012"}, l.Annotations)
}
