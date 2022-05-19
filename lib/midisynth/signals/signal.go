package signals

import "gitlab.com/avoronkov/waver/lib/notes"

type Signal struct {
	Instrument   int
	Note         notes.Note
	Sample       string `json:",omitempty"`
	DurationBits int
	Amp          float64
	// Manual control section
	Manual bool `json:",omitempty"`
	Stop   bool `json:",omitempty"`
}
