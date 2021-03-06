package waves

type NoteCtx struct {
	Freq   float64
	Period float64
	Amp    float64
	Dur    float64

	Channel int
	AbsTime float64
}

func NewNoteCtx(freq float64, amp float64, dur float64, absTime float64) *NoteCtx {
	return &NoteCtx{
		Freq:    freq,
		Period:  1.0 / freq,
		Amp:     amp,
		Dur:     dur,
		AbsTime: absTime,
	}
}

func (c *NoteCtx) SetFrequency(freq float64) {
	c.Freq = freq
	c.Period = 1.0 / freq
}
