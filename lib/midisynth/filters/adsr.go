package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type AdsrFilter struct {
	AttackLevel float64
	DecayLevel  float64

	AttackLen  float64
	DecayLen   float64
	SusteinLen float64
	ReleaseLen float64
}

var _ FilterAdsr = (*AdsrFilter)(nil)

func NewAdsrFilter(opts ...func(*AdsrFilter)) Filter {
	f := &AdsrFilter{
		AttackLevel: 1.0,
		DecayLevel:  1.0,
		ReleaseLen:  1.0,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (af *AdsrFilter) Apply(w waves.Wave) waves.Wave {
	return &adsrImpl{
		wave: w,
		opts: af,
	}
}

func (af *AdsrFilter) IsAdsr() {}

type adsrImpl struct {
	wave waves.Wave
	opts *AdsrFilter
}

func (i *adsrImpl) Value(tm float64, ctx *waves.NoteCtx) float64 {
	amp := 0.0
	dur := ctx.Dur
	o := i.opts

	adsrLen := o.AttackLen + o.DecayLen + o.SusteinLen + o.ReleaseLen

	if attackLen := o.AttackLen * dur / adsrLen; tm >= 0 && tm < attackLen {
		// attack
		amp = tm * i.opts.AttackLevel / attackLen
	} else if tm < (o.AttackLen+o.DecayLen)*dur/adsrLen {
		// decay
		amp = o.AttackLevel - (o.AttackLevel-o.DecayLevel)*(tm-(o.AttackLen*dur)/adsrLen)/(o.DecayLen*dur/adsrLen)
	} else if tm < (o.AttackLen+o.DecayLen+o.SusteinLen)*dur/adsrLen {
		// sustein
		amp = o.DecayLevel
	} else if tm < dur {
		// release
		amp = (dur - tm) * o.DecayLevel * adsrLen / (dur * o.ReleaseLen)
	} else {
		return math.NaN()
	}

	return i.wave.Value(tm, ctx) * amp * ctx.Amp
}

// Options
func AdsrAttackLevel(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.AttackLevel = v
	}
}

func AdsrDecayLevel(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.DecayLevel = v
	}
}

func AdsrAttackLen(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.AttackLen = v
	}
}

func AdsrDecayLen(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.DecayLen = v
	}
}

func AdsrSusteinLen(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.SusteinLen = v
	}
}

func AdsrReleaseLen(v float64) func(*AdsrFilter) {
	return func(f *AdsrFilter) {
		f.ReleaseLen = v
	}
}
