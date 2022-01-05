package waves

import "math"

type sine struct {
	hz     float64
	period float64
}

func Sine(hz float64) Wave {
	return &sine{
		hz:     hz,
		period: 1.0 / hz,
	}
}

func (s *sine) Value(t float64) float64 {
	x := 2.0 * math.Pi * t / s.period
	return math.Sin(x)
}
