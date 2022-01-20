package waves

import "math"

var Square = WaveFn(func(t float64, ctx *NoteCtx) float64 {
	y := t / ctx.Period
	if y-math.Floor(y) < 0.5 {
		return 1.0
	}
	return -1.0
})
