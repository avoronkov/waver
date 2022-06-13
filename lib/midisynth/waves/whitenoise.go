package waves

import (
	"math/rand"
	"time"
)

type whiteNoiseImpl struct {
	rng *rand.Rand
}

func NewWhiteNoise() Wave {
	return &whiteNoiseImpl{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (n *whiteNoiseImpl) Value(t float64, ctx *NoteCtx) float64 {
	sign := 1.0
	if n.rng.Intn(2) == 0 {
		sign = -1.0
	}
	return n.rng.Float64() * sign
}

var WhiteNoise = NewWhiteNoise()
