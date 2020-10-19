package logbus

import (
	"context"
	"fmt"
)

var NoInputData = fmt.Errorf("no input data")

type Input interface {
	Read(ctx context.Context) (*StandardLog, error)
	Close() error
}
