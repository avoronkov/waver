package seq

func WithTempo(n int) func(*Sequencer) {
	return func(s *Sequencer) {
		s.tempo = n
	}
}

func WithStart(bit int64) func(*Sequencer) {
	return func(s *Sequencer) {
		s.bit = bit
	}
}
