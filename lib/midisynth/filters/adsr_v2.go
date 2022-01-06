package filters

import "waver/lib/midisynth/waves/v2"

type AdsrFilter struct {
	AttackLevel float64
	DecayLevel  float64

	AttackLen  float64
	DecayLen   float64
	SusteinLen float64
	ReleaseLen float64
}

func NewAdsrFilter(opts ...func(*AdsrFilter)) Filter {
	f := &AdsrFilter{
		AttackLevel: 1.0,
		DecayLevel:  1.0,
		ReleaseLen:  1.0,
	}

	// TODO handle opts

	return f
}

func (af *AdsrFilter) Apply(w waves.Wave) waves.Wave {
	return &adsrImpl{
		wave: w,
		opts: af,
	}
}

type adsrImpl struct {
	wave waves.Wave
	opts *AdsrFilter
}

func (i *adsrImpl) Value(tm float64, ctx *waves.NoteCtx) float64 {
	amp := 0.0
	dur := ctx.Dur
	o := i.opts

	if attackLen := o.AttackLen * dur; tm < attackLen {
		// attack
		amp = tm * i.opts.AttackLevel / attackLen
	} else if tm < (o.AttackLen+o.DecayLen)*dur {
		// decay
		amp = o.AttackLevel - (o.AttackLevel-o.DecayLevel)*(tm-(o.AttackLen*dur))/(o.DecayLen*dur)
	} else if tm < (o.AttackLen+o.DecayLen+o.SusteinLen)*dur {
		// sustein
		amp = o.DecayLevel
	} else if tm < (o.AttackLen+o.DecayLen+o.SusteinLen+o.ReleaseLen)*dur {
		// release
		j := tm - (o.AttackLen-o.DecayLen-o.SusteinLen)*dur
		amp = o.DecayLevel - o.DecayLevel*j/(o.ReleaseLen*dur)
	}

	return i.wave.Value(tm, ctx) * amp * ctx.Amp
}
