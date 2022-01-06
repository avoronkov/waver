package filters

import "waver/lib/midisynth/waves/v2"

type Filter interface {
	Apply(w waves.Wave) waves.Wave
}
