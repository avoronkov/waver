package waves

type Wave interface {
	Value(tm float64, ctx *NoteCtx) float64
}

type WithDuration interface {
	// Duration in seconds
	Duration(ctx *NoteCtx) float64
}

type WaveControlled interface {
	Wave
	Release()
}
