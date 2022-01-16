package midi

func Edo19() func(p *Proc) {
	return func(p *Proc) {
		p.keyMap = KeyMap19
	}
}
