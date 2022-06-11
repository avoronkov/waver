package filters

import (
	"fmt"
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
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

func (f DelayFilter) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o []func(*DelayFilter)
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "int", "interval":
				v := float64Of(value)
				o = append(o, DelayInterval(v))
			case "times":
				v := value.(int)
				o = append(o, DelayTimes(v))
			case "fade":
				v := float64Of(value)
				o = append(o, DelayFadeOut(v))
			default:
				return nil, fmt.Errorf("Unknown Delay parameter: %v", param)
			}
		}
	}
	return NewDelayFilter(o...), nil
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
