package midi

func Edo19() func(p *Proc) {
	return func(p *Proc) {
		p.keyMap = KeyMap19
	}
}

func WithDispatcher(d Dispatcher) func(p *Proc) {
	return func(p *Proc) {
		p.dispatcher = d
	}
}

func WithChannel(ch chan<- string) func(p *Proc) {
	return func(p *Proc) {
		p.ch = ch
	}
}
