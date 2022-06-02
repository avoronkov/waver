package parser

func WithFileInput(file string) func(*Parser) {
	return func(p *Parser) {
		p.file = file
	}
}
