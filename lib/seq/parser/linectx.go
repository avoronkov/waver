package parser

type LineCtx struct {
	Num       int
	Fields    []string
	GlobalCtx map[string]interface{}
}

func (c *LineCtx) Shift(n int) *LineCtx {
	return &LineCtx{
		Num:       c.Num,
		Fields:    c.Fields[n:],
		GlobalCtx: c.GlobalCtx,
	}
}

func (c *LineCtx) Len() int {
	return len(c.Fields)
}
