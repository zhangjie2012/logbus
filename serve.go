package logbus

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	readSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logbus_read_success_count",
		Help: "The logbus read log success count",
	})
	readIdle = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logbus_read_idle_loop_count",
		Help: "The logbus read idle loop count",
	})
	readFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logbus_read_fail_count",
		Help: "The logbus read log failure count",
	})

	writeSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logbus_write_success_count",
		Help: "The logbus write log success count",
	})
	writeFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logbus_write_fail_count",
		Help: "The logbus write log failue count",
	})
)

func tickRead(err error) {
	if err != nil {
		if err != NoInputData {
			readFail.Inc()
		} else {
			readIdle.Inc()
		}
	} else {
		readSuccess.Inc()
	}
}

func tickWrite(err error) {
	if err != nil {
		writeFail.Inc()
	} else {
		writeSuccess.Inc()
	}
}

func Serve(ctx context.Context, in Input, outputs []Output) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			l, err := readOne(in)
			tickRead(err)
			if err != nil {
				if err != NoInputData {
					log.Print("read", err)
				}
				break
			}

			for _, o := range outputs {
				err := o.Write(l)
				tickWrite(err)
				if err != nil {
					log.Print("write", err)
				}
			}
		}
	}
}

func readOne(in Input) (*StdLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return in.Read(ctx)
}
