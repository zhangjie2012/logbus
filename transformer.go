package logbus

type TransformerFunc func(l *StdLog) (*StdLog, bool)

func DefaultTransformer(l *StdLog) (*StdLog, bool) {
	return l, true
}

func StatLogTransformer(l *StdLog) (*StdLog, bool) {
	if l.StateId != StateIdInvalid {
		return l, true
	}
	return nil, false
}
