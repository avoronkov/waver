package waves

import "math"

var SineSine = WaveFn(func(t float64, ctx *NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	x1 := (math.Sin(x) + 1.0) * 2.0 * math.Pi
	x2 := (math.Sin(x1) + 1.0) * 2.0 * math.Pi
	return math.Sin(x2)
})
