package pragma

func WithTempoSetter(ts TempoSetter) func(*PragmaParser) {
	return func(pp *PragmaParser) {
		pp.tempoSetters = append(pp.tempoSetters, ts)
	}
}

func WithInstrumentSet(is InstrumentSet) func(*PragmaParser) {
	return func(pp *PragmaParser) {
		pp.instSet = is
	}
}
