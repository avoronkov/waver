package signals

type Signal struct {
	Instrument   int
	Octave       int
	Note         string
	Sample       string
	DurationBits int
	Amp          float64
	// Manual control section
	Manual bool
	Stop   bool
}
