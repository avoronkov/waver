package waves

import "math"

var Sine = WaveFn(func(t float64, ctx *NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	return math.Sin(x)
})
