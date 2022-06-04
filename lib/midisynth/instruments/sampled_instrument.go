package instruments

import "github.com/avoronkov/waver/lib/midisynth/waves"

type SampledInstrument struct {
}

var _ Interface = (*SampledInstrument)(nil)

func (s *SampledInstrument) Wave() waves.Wave {
	return nil
}

func (s *SampledInstrument) WaveControlled() waves.WaveControlled {
	return nil
}
