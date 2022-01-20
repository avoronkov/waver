package waves

import "math"

var SemiSine = WaveFn(func(t float64, ctx *NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	return 2.0*math.Abs(math.Sin(x)) - 1.0
})
