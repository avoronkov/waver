package pragma

func WithTempoSetter(ts TempoSetter) func(*PragmaParser) {
	return func(pp *PragmaParser) {
		pp.tempoSetters = append(pp.tempoSetters, ts)
	}
}
