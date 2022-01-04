package filters

import (
	"waver/lib/midisynth/waves"
)

type Adsr struct {
	wave waves.Wave
	opts *AdsrOpts
}

type AdsrOpts struct {
	AttackLevel float64
	DecayLevel  float64

	AttackLen  float64
	DecayLen   float64
	SusteinLen float64
	ReleaseLen float64
}

func NewAdsr(wave waves.Wave, o ...func(*AdsrOpts)) waves.WaveDuration {
	opts := &AdsrOpts{
		AttackLevel: 1.0,
		DecayLevel:  1.0,
	}
	for _, opt := range o {
		opt(opts)
	}
	return &Adsr{
		wave: wave,
		opts: opts,
	}
}

func (a *Adsr) Value(t float64) float64 {
	amp := 0.0
	o := a.opts
	if t < a.opts.AttackLen {
		// attack
		amp = t * o.AttackLevel / o.AttackLen
	} else if t < o.AttackLen+o.DecayLen {
		// decay
		amp = o.AttackLevel - (o.AttackLevel-o.DecayLevel)*(t-o.AttackLen)/o.DecayLen
	} else if t < o.AttackLen+o.DecayLen+o.SusteinLen {
		// sustein
		amp = o.DecayLevel
	} else if t < a.Duration() {
		j := t - o.AttackLen - o.DecayLen - o.SusteinLen
		amp = o.DecayLevel - o.DecayLevel*j/o.ReleaseLen
	}

	return a.wave.Value(t) * amp
}

func (a *Adsr) Duration() float64 {
	return a.opts.AttackLen + a.opts.DecayLen + a.opts.SusteinLen + a.opts.ReleaseLen
}

// ADSR options
type adsrOpt func(*AdsrOpts)

func AdsrAttackLevel(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.AttackLevel = v
	}
}

func AdsrDecayLevel(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.DecayLevel = v
	}
}

func AdsrAttackLen(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.AttackLen = v
	}
}

func AdsrDecayLen(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.DecayLen = v
	}
}
func AdsrSusteinLen(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.SusteinLen = v
	}
}
func AdsrReleaseLen(v float64) func(*AdsrOpts) {
	return func(o *AdsrOpts) {
		o.ReleaseLen = v
	}
}
