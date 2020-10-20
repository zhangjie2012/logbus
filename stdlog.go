package logbus

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// only state log has one valid id
const StateIdInvalid string = "_invalid"

// StandardLog standard log format
//   - AppName define where log from
//   - Annotations for struct log or message expand
type StdLog struct {
	AppName     string                 `msgpack:"appname" bson:"appname"`
	Timestamp   int64                  `msgpack:"timestamp" bson:"timestamp"`
	Level       string                 `msgpack:"level" bson:"level"`
	StateId     string                 `msgpack:"stateid" bson:"stateid"`
	Caller      string                 `msgpack:"caller" bson:"caller"`
	Message     string                 `msgpack:"message" bson:"message"`
	Annotations map[string]interface{} `msgpack:"annotations" bson:"annotations"`
}

func newStdLog(appName string, t time.Time, metadata logrus.Fields, caller *runtime.Frame, level logrus.Level, message string) *StdLog {
	caller_ := ""
	if caller != nil {
		caller_ = fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)
	}

	stateId, ok := metadata["stateid"].(string)
	if !ok {
		stateId = StateIdInvalid
	}

	annotations := map[string]interface{}{}
	for k, v := range metadata {
		if k != "stateid" {
			annotations[k] = v
		}
	}

	l := StdLog{
		AppName:     appName,
		Timestamp:   t.UnixNano(),
		Level:       level.String(),
		StateId:     stateId,
		Caller:      caller_,
		Message:     message,
		Annotations: annotations,
	}

	return &l
}

// StdLogWash redefined "LogWashFunc" in logrusredis-hook
// use "StandardLog" replace logrusredis-hook's "DefaultLogS"
func StdLogWash(appName string, t time.Time, metadata logrus.Fields, caller *runtime.Frame, level logrus.Level, message string) []byte {
	l := newStdLog(appName, t, metadata, caller, level, message)
	bs, err := msgpack.Marshal(&l)
	if err != nil {
		return nil
	}
	return bs
}
