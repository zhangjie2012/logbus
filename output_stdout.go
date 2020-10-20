package logbus

import (
	"fmt"
	"time"
)

type StdoutOutput struct {
	transformer TransformerFunc
}

func NewStdoutOutput(transformer TransformerFunc) (Output, error) {
	return &StdoutOutput{transformer: transformer}, nil
}

func (out *StdoutOutput) Write(l *StdLog) error {
	lb, ok := out.transformer(l)
	if ok {
		t := time.Unix(lb.Timestamp/1e9, lb.Timestamp%1e9).Format(time.RFC3339)
		fmt.Printf("appname=%s, timestamp=%s, caller=%s, level=%s, id=%s, annotations=%s, message=%s\n",
			lb.AppName, t, lb.Caller, lb.Level, lb.StateId, lb.Annotations, lb.Message)
	}
	return nil
}

func (out *StdoutOutput) Close() error {
	return nil
}
