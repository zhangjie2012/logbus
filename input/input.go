package input

import (
	"context"
	"fmt"

	"github.com/zhangjie2012/logbus"
)

var NoInputData = fmt.Errorf("no input data")

type Input interface {
	Read(ctx context.Context) (*logbus.StandardLog, error)
	Close() error
}
