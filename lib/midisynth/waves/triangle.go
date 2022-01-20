package waves

import "math"

var Triangle = WaveFn(func(t float64, ctx *NoteCtx) (val float64) {
	// Start from 0-point.
	y := (t / ctx.Period) + 0.25
	y = y - math.Floor(y)
	if y < 0.5 {
		val = -1.0 + 4.0*y
	} else {
		val = 1.0 - 4.0*(y-0.5)
	}
	return
})
