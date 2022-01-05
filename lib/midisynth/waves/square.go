package waves

import "math"

type square struct {
	hz     float64
	period float64
}

func Square(hz float64) Wave {
	return &square{
		hz:     hz,
		period: 1.0 / hz,
	}
}

func (s *square) Value(t float64) float64 {
	y := t / s.period
	if y-math.Floor(y) < 0.5 {
		return 1.0
	}
	return -1.0
}
