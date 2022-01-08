package filters

import "gitlab.com/avoronkov/waver/lib/midisynth/waves"

// Ring (amplitude) modulation.
type Ring struct {
	Carrier    waves.Wave
	CarrierCtx *waves.NoteCtx
}

func NewRing(carrier waves.Wave, freq float64) Filter {
	return &Ring{
		Carrier:    carrier,
		CarrierCtx: waves.NewNoteCtx(freq, 1.0, -1.0),
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
	return i.input.Value(t, ctx) * i.opts.Carrier.Value(t, i.opts.CarrierCtx)
}
