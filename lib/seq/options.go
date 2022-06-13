package seq

import "github.com/avoronkov/waver/lib/midisynth/signals"

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

func WithChannel(ch chan<- signals.Interface) func(*Sequencer) {
	return func(s *Sequencer) {
		s.ch = ch
	}
}
