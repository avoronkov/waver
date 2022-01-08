package filters

import "waver/lib/midisynth/waves"

type Filter interface {
	Apply(w waves.Wave) waves.Wave
}
