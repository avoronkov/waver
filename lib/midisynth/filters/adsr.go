package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type AdsrFilter struct {
	AttackLevel float64 `option:"attackLevel"`
	DecayLevel  float64 `option:"decayLevel"`

	AttackLen  float64 `option:"attackLen"`
	DecayLen   float64 `option:"decayLen"`
	SustainLen float64 `option:"sustainLen"`
	ReleaseLen float64 `option:"releaseLen"`
}

var _ FilterAdsr = (*AdsrFilter)(nil)

func (AdsrFilter) New() Filter {
	return &AdsrFilter{
		AttackLevel: 1.0,
		DecayLevel:  1.0,
		ReleaseLen:  1.0,
	}
}

func (AdsrFilter) Desc() string {
	return `The most common envelope generator controlled with four parameters: attack, decay, sustain and release (ADSR).`
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

	adsrLen := o.AttackLen + o.DecayLen + o.SustainLen + o.ReleaseLen

	if attackLen := o.AttackLen * dur / adsrLen; tm >= 0 && tm < attackLen {
		// attack
		amp = tm * i.opts.AttackLevel / attackLen
	} else if tm < (o.AttackLen+o.DecayLen)*dur/adsrLen {
		// decay
		amp = o.AttackLevel - (o.AttackLevel-o.DecayLevel)*(tm-(o.AttackLen*dur)/adsrLen)/(o.DecayLen*dur/adsrLen)
	} else if tm < (o.AttackLen+o.DecayLen+o.SustainLen)*dur/adsrLen {
		// sustain
		amp = o.DecayLevel
	} else if tm < dur {
		// release
		amp = (dur - tm) * o.DecayLevel * adsrLen / (dur * o.ReleaseLen)
	} else {
		return math.NaN()
	}

	return i.wave.Value(tm, ctx) * amp
}
