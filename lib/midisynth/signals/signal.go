package signals

type Signal struct {
	Instrument   int
	Octave       int
	Note         string
	Sample       string
	DurationBits int
	Amp          float64
	Stop         bool
}
