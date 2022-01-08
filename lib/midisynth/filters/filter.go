package filters

import "gitlab.com/avoronkov/waver/lib/midisynth/waves"

type Filter interface {
	Apply(w waves.Wave) waves.Wave
}
