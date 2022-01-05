package waves

import "math"

type sawtooth struct {
	period float64
	invert float64
}

func Saw(hz float64, invert bool) Wave {
	s := &sawtooth{
		period: 1.0 / hz,
		invert: 1.0,
	}
	if invert {
		s.invert = -1.0
	}
	return s
}

func (s *sawtooth) Value(t float64) float64 {
	y := 0.0
	if s.invert > 0.0 {
		y = (t / s.period) + 0.25
	} else {
		y = t / s.period
	}
	y = y - math.Floor(y)
	return (y*2.0 - 1.0) * s.invert
}
