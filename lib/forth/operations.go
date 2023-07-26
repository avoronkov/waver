package forth

import (
	"fmt"
)

func Number(n int) func(*Forth) error {
	return func(f *Forth) error {
		// log.Printf("Number %v: %v", n, s.stack)
		f.Stack.Push(n)
		return nil
	}
}

func Dup(f *Forth) error {
	// log.Printf("Dup: %v", s.stack)
	v, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	f.Stack.Push(v)
	f.Stack.Push(v)
	return nil
}

func Drop(f *Forth) error {
	// log.Printf("Drop: %v", s.stack)
	_, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	return nil
}

func Loop(fns []StackFn) func(*Forth) error {
	return func(f *Forth) error {
		// log.Printf("Loop: %v", s.stack)
		for {
			// check value on top of the stack
			v, ok := f.Stack.Pop()
			if !ok {
				return EmptyStack
			}
			if v == 0 {
				break
			}

			// execute loop actions
			for _, fn := range fns {
				if err := fn(f); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func Sequence(fns []StackFn) func(*Forth) error {
	return func(f *Forth) error {
		for _, fn := range fns {
			if err := fn(f); err != nil {
				return err
			}
		}
		return nil
	}
}

func Message(msg string) func(*Forth) error {
	return func(f *Forth) error {
		fmt.Print(msg)
		return nil
	}
}

func ShowTop(f *Forth) error {
	v, ok := f.Stack.Top()
	if !ok {
		return EmptyStack
	}
	fmt.Println(v)
	return nil
}

func ShowStack(f *Forth) error {
	fmt.Println(f.Stack.stack)
	return nil
}

func Plus(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	b, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	f.Stack.Push(a + b)
	return nil
}

func Minus(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	b, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	f.Stack.Push(b - a)
	return nil
}

func Multiply(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	b, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	f.Stack.Push(a * b)
	return nil
}

func And(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	b, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	res := 0
	if a != 0 && b != 0 {
		res = 1
	}
	f.Stack.Push(res)
	return nil
}

func Or(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	b, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	res := 0
	if a != 0 || b != 0 {
		res = 1
	}
	f.Stack.Push(res)
	return nil
}

func Not(f *Forth) error {
	a, ok := f.Stack.Pop()
	if !ok {
		return EmptyStack
	}
	res := 0
	if a == 0 {
		res = 1
	}
	f.Stack.Push(res)
	return nil
}

func Function(name string) func(*Forth) error {
	return func(f *Forth) error {
		fn, ok := f.funcs[name]
		if !ok {
			return UnknownFunction(name)
		}
		return fn(f)
	}
}
