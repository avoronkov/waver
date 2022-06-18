package instruments

import (
	"log"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Set struct {
	waved map[string]*Instrument
}

func NewSet() *Set {
	return &Set{
		waved: make(map[string]*Instrument),
	}
}

func (s *Set) Wave(inst string) (waves.Wave, bool) {
	in, ok := s.waved[inst]
	if !ok {
		return nil, false
	}
	return in.Wave(), true
}

func (s *Set) WaveControlled(inst string) (waves.WaveControlled, bool) {
	in, ok := s.waved[inst]
	if !ok {
		return nil, false
	}
	return in.WaveControlled(), true
}

func (s *Set) AddInstrument(n string, inst *Instrument) {
	log.Printf("Set: AddInstrument %v", n)
	s.waved[n] = inst
}
