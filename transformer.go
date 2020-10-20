package logbus

type TransformerFunc func(l *StandardLog) (*StandardLog, bool)

func DefaultTransformer(l *StandardLog) (*StandardLog, bool) {
	return l, true
}

func StatLogTransformer(l *StandardLog) (*StandardLog, bool) {
	if l.StateId != StateIdInvalid {
		return l, true
	}
	return nil, false
}
