package parser

type LineCtx struct {
	Num    int
	Fields []string
}

func (c *LineCtx) Shift(n int) *LineCtx {
	return &LineCtx{
		Num:    c.Num,
		Fields: c.Fields[n:],
	}
}

func (c *LineCtx) Len() int {
	return len(c.Fields)
}
