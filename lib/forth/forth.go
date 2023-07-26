package forth

type StackFn func(s *Forth) error

type Forth struct {
	Stack *Stack

	funcs map[string]StackFn

	program []StackFn
}

func New(opts ...func(*Forth)) *Forth {
	f := &Forth{
		Stack: NewStack(),
		funcs: map[string]StackFn{},
	}

	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithFuncs(funcs map[string]StackFn) func(*Forth) {
	return func(f *Forth) {
		for name, fn := range funcs {
			f.funcs[name] = fn
		}
	}
}

func WithProgram(program []StackFn) func(*Forth) {
	return func(f *Forth) {
		f.program = program
	}
}

func (f *Forth) Run() error {
	for _, fn := range f.program {
		if err := fn(f); err != nil {
			return err
		}
	}
	return nil
}
