package seq

func WithTempo(n int) func(*Sequencer) {
	return func(s *Sequencer) {
		s.tempo = n
	}
}
