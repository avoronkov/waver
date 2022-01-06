package waves

type Wave interface {
	Value(tm float64, ctx *NoteCtx) float64
}
