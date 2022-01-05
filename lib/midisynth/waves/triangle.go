package waves

import "math"

type triangle struct {
	hz     float64
	period float64
}

func Triangle(hz float64) Wave {
	return &triangle{
		hz:     hz,
		period: 1.0 / hz,
	}
}

func (w *triangle) Value(t float64) (val float64) {
	// Start from 0-point.
	y := (t / w.period) + 0.25
	y = y - math.Floor(y)
	if y < 0.5 {
		val = -1.0 + 4.0*y
	} else {
		val = 1.0 - 4.0*(y-0.5)
	}
	return
}
