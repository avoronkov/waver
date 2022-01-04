package envelops

import (
	"io"
	"waver/lib/midisynth/waves"
)

type Envelop interface {
	Wrap(wave waves.Wave) io.Reader
}
