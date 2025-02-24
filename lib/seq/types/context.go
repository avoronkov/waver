package types

import (
	"fmt"
	"maps"
)

func NewContext(opts ...func(*contextImpl)) Context {
	impl := &contextImpl{
		data:   map[string]ValueFn{},
		values: map[string]Value{},
	}
	for _, opt := range opts {
		opt(impl)
	}
	if impl.globalCtx == nil {
		impl.globalCtx = make(map[string]any)
	}
	return impl
}

func WithGlobalContext(ctx map[string]any) func(*contextImpl) {
	return func(impl *contextImpl) {
		impl.globalCtx = ctx
	}
}

type contextImpl struct {
	data      map[string]ValueFn
	values    map[string]Value
	globalCtx map[string]any
}

func (c *contextImpl) Put(name string, fn ValueFn) error {
	/*
		if _, exists := c.data[name]; exists {
			return fmt.Errorf("Cannot re-assign variable: %v", name)
		}
	*/
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

func (c *contextImpl) GlobalPut(name string, value any) {
	c.globalCtx[name] = value
}

func (c *contextImpl) GlobalGet(name string) (value any, ok bool) {
	value, ok = c.globalCtx[name]
	return
}

func (c *contextImpl) Copy() Context {
	cc := &contextImpl{
		data:      map[string]ValueFn{},
		values:    map[string]Value{},
		globalCtx: c.globalCtx,
	}
	maps.Copy(cc.data, c.data)
	maps.Copy(cc.values, c.values)
	return cc
}
