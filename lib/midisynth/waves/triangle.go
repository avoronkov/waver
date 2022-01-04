package waves

import "math"

type Triangle struct {
	hz     float64
	period float64
}

func NewTriangle(hz float64) Wave {
	return &Triangle{
		hz:     hz,
		period: 1.0 / hz,
	}
}

func (w *Triangle) Value(t float64) (val float64) {
	// Start from 0-point.
	y := ((t / w.period) + 0.25) / 1.0
	y = y - math.Floor(y)
	if y < 0.5 {
		val = -1.0 + 4.0*y
	} else {
		val = 1.0 - 4.0*(y-0.5)
	}
	return
}
