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

type WaveFn func(t float64, ctx *NoteCtx) float64

func (f WaveFn) Value(t float64, ctx *NoteCtx) float64 {
	return f(t, ctx)
}
