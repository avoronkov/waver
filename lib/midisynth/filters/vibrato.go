package filters

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type VibratoFilter struct {
	Shifter    waves.Wave
	ShifterCtx *waves.NoteCtx
}

func NewVibrato(w waves.Wave, freq, amp float64) Filter {
	return &VibratoFilter{
		Shifter:    w,
		ShifterCtx: waves.NewNoteCtx(freq, amp, -1.0),
	}
}

func (v *VibratoFilter) Apply(input waves.Wave) waves.Wave {
	return &vibratoImpl{
		input: input,
		opts:  v,
	}
}

type vibratoImpl struct {
	input waves.Wave
	opts  *VibratoFilter
}

func (i *vibratoImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	newFreq := (i.opts.Shifter.Value(t, i.opts.ShifterCtx)*i.opts.ShifterCtx.Amp + 1.0) * ctx.Freq
	newCtx := waves.NewNoteCtx(newFreq, ctx.Amp, ctx.Dur)
	res := i.input.Value(t, newCtx)
	// log.Printf("newFreq(t=%v) = %v (%v) -> %v", t, newFreq, ctx.Freq, res)
	return res
}
