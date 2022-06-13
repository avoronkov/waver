package parser

func WithFileInput(file string) func(*Parser) {
	return func(p *Parser) {
		p.file = file
	}
}

func WithTempoSetter(ts TempoSetter) func(*Parser) {
	return func(pp *Parser) {
		pp.tempoSetters = append(pp.tempoSetters, ts)
	}
}

func WithInstrumentSet(is InstrumentSet) func(*Parser) {
	return func(pp *Parser) {
		pp.instSet = is
	}
}
