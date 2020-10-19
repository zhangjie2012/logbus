package logbus

type Output interface {
	Write(l *StandardLog) error
	Close() error
}
