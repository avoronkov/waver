package seq

func WithTempo(n int) func(*Sequencer) {
	return func(s *Sequencer) {
		s.tempo = n
	}
}

func WithStart(bit int64) func(*Sequencer) {
	return func(s *Sequencer) {
		s.startingBit = bit
	}
}

func WithShowingBits(n int64) func(*Sequencer) {
	return func(s *Sequencer) {
		s.showBits = n
	}
}
