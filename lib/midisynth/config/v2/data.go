package config

/*

instruments:
  1:
    wave:
	  sine:
	filters:
	- adsr:
	    releaseLen: 1.0
	- delay:
		interval: 0.5
		times: 4
		fadeOut: 0.5

samples:
  4k:
	sample: "samples/4-kick.wav"
  4s:
    sample: "samples/4-snare.wav"
  4h:
    sample: "samples/4-hat.wav"
*/

type Data struct {
	Instruments map[string]Instrument `yaml:"instruments"`
	Samples     map[string]SampleData `yaml:"samples"`
}

type Instrument struct {
	Wave    string   `yaml:"wave"`
	Filters []Filter `yaml:"filters"`
}

type SampleData struct {
	Sample  string   `yaml:"sample"`
	Filters []Filter `yaml:"filters"`
}

type Filter map[string]map[string]interface{}
