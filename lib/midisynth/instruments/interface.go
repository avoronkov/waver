package instruments

import "github.com/avoronkov/waver/lib/midisynth/waves"

type Interface interface {
	Wave() waves.Wave
	WaveControlled() waves.WaveControlled
}
