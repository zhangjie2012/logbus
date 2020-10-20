package logbus

import (
	"fmt"
	"time"
)

type StdoutOutput struct{}

func NewStdoutOutput() (Output, error) {
	return &StdoutOutput{}, nil
}

func (out *StdoutOutput) Write(l *StdLog) error {
	t := time.Unix(l.Timestamp/1e9, l.Timestamp%1e9)
	fmt.Printf("timestamp=%s, level=%s, message=%s\n", t, l.Level, l.Message)
	return nil
}

func (out *StdoutOutput) Close() error {
	return nil
}
