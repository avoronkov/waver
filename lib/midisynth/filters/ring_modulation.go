package filters

import "github.com/avoronkov/waver/lib/midisynth/waves"

// Ring (amplitude) modulation.
type Ring struct {
	Carrier    waves.Wave
	CarrierCtx *waves.NoteCtx
}

func NewRing(carrier waves.Wave, freq, amp float64) Filter {
	return &Ring{
		Carrier:    carrier,
		CarrierCtx: waves.NewNoteCtx(freq, amp, -1.0),
	}
}

func (rf *Ring) Apply(input waves.Wave) waves.Wave {
	return &ringImpl{
		input: input,
		opts:  rf,
	}
}

type ringImpl struct {
	input waves.Wave
	opts  *Ring
}

func (i *ringImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	mp := (2.0 - i.opts.CarrierCtx.Amp + i.opts.CarrierCtx.Amp*i.opts.Carrier.Value(t, i.opts.CarrierCtx)) / 2.0
	return i.input.Value(t, ctx) * mp
}
