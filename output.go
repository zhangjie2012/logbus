package logbus

type Output interface {
	Write(l *StdLog) error
	Close() error
}
