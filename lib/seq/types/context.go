package types

import "fmt"

func NewContext() Context {
	return &contextImpl{
		data:   map[string]ValueFn{},
		values: map[string]Value{},
	}
}

type contextImpl struct {
	data   map[string]ValueFn
	values map[string]Value
}

func (c *contextImpl) Put(name string, fn ValueFn) error {
	if _, exists := c.data[name]; exists {
		return fmt.Errorf("Cannot re-assign variable: %v", name)
	}
	c.data[name] = fn
	return nil
}

func (c *contextImpl) Get(name string, bit int64) Value {
	if val, ok := c.values[name]; ok {
		return val
	}
	fn, ok := c.data[name]
	if !ok {
		panic(fmt.Errorf("Unknown function: %v", name))
	}
	value := fn.Val(bit, c)
	c.values[name] = value
	return value
}
