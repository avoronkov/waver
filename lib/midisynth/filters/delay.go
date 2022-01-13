package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type DelayFilter struct {
	// In seconds
	Interval float64

	Times int

	// [0.0..1.0]
	FadeOut float64
}

func NewDelayFilter(opts ...func(*DelayFilter)) Filter {
	f := &DelayFilter{}
	for _, o := range opts {
		o(f)
	}
	return f
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
var _ waves.WithDuration = (*delayImpl)(nil)

func (d *delayImpl) Value(tm float64, ctx *waves.NoteCtx) float64 {
	value := d.wave.Value(tm, ctx)
	if math.IsNaN(value) && tm < float64(d.opts.Times)*d.opts.Interval {
		value = 0
	}

	multiplier := 1.0
	t := tm
	for i := 0; i < d.opts.Times; i++ {
		multiplier *= d.opts.FadeOut
		t -= d.opts.Interval
		if t < 0.0 {
			break
		}
		v := d.wave.Value(t, ctx)
		if math.IsNaN(v) {
			continue
		}
		value += v * multiplier
	}

	return value
}

func (d *delayImpl) Duration(ctx *waves.NoteCtx) float64 {
	// TODO better error handling
	origDuration := 0.0
	if withDur, ok := d.wave.(waves.WithDuration); ok {
		origDuration = withDur.Duration(ctx)
	} else {
		origDuration = ctx.Dur
	}
	return float64(d.opts.Times)*d.opts.Interval + origDuration
}

// Options

func DelayInterval(d float64) func(*DelayFilter) {
	return func(f *DelayFilter) {
		f.Interval = d
	}
}

func DelayTimes(n int) func(*DelayFilter) {
	return func(f *DelayFilter) {
		f.Times = n
	}
}

func DelayFadeOut(out float64) func(*DelayFilter) {
	return func(f *DelayFilter) {
		f.FadeOut = out
	}
}
