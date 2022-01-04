package waves

import "math"

type SineWave struct {
	hz float64
	// amp [0, 1]
	amp float64

	period float64
}

func NewSineWave(hz float64, amp float64) Wave {
	return &SineWave{
		hz:     hz,
		amp:    amp,
		period: 1.0 / hz,
	}
}

func (s *SineWave) Value(t float64) float64 {
	x := 2.0 * math.Pi * t / s.period
	return s.amp * math.Sin(x)
}
