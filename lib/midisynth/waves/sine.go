package waves

import "math"

type Sine struct{}

func (s *Sine) Value(t float64, ctx *NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	return math.Sin(x)
}
