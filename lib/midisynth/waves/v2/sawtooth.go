package waves

import "math"

type Saw struct{}

func (s *Saw) Value(t float64, ctx *NoteCtx) float64 {
	y := (t / ctx.Period) + 0.25
	y = y - math.Floor(y)
	return y*2.0 - 1.0
}
