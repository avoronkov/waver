package waves

import "math"

type Sine struct {
	hz     float64
	period float64
}

func NewSine(hz float64) Wave {
	return &Sine{
		hz:     hz,
		period: 1.0 / hz,
	}
}

func (s *Sine) Value(t float64) float64 {
	x := 2.0 * math.Pi * t / s.period
	return math.Sin(x)
}
