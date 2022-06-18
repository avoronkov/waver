package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type DelayFilter struct {
	// In seconds
	Interval float64 `option:"interval,int"`

	Times int `option:"times"`

	// [0.0..1.0]
	FadeOut float64 `option:"fade,feedback"`
}

func (DelayFilter) New() Filter {
	return &DelayFilter{
		Interval: 1.0,
		Times:    1,
		FadeOut:  0.5,
	}
}

func (df *DelayFilter) Apply(w waves.Wave) waves.Wave {
	return &delayImpl{
		wave: w,
		opts: df,
	}
}

type delayImpl struct {
	wave waves.Wave
	opts *DelayFilter
}

var _ waves.Wave = (*delayImpl)(nil)

func (d *delayImpl) Value(tm float64, ctx *waves.NoteCtx) (res float64) {
	value := d.wave.Value(tm, ctx)
	has := false
	if math.IsNaN(value) {
		value = 0
		if tm < float64(d.opts.Times+1)*d.opts.Interval {
			has = true
		}
	} else {
		has = true
	}

	multiplier := 1.0
	t := tm
	for i := 1; i <= d.opts.Times; i++ {
		t -= d.opts.Interval
		multiplier *= d.opts.FadeOut
		if t < 0.0 {
			break
		}
		v := d.wave.Value(t, ctx)
		if math.IsNaN(v) {
			continue
		}
		value += v * multiplier
		has = true
	}

	if !has {
		return math.NaN()
	}

	return value
}
