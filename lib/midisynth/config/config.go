package config

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
	"gitlab.com/avoronkov/waver/lib/watch"
	"gitlab.com/avoronkov/waver/static"

	yaml "gopkg.in/yaml.v3"
)

type InstrumentSet interface {
	AddInstrument(n int, in *instruments.Instrument)
	AddSampledInstrument(name string, in *instruments.Instrument)
}

type Config struct {
	m         InstrumentSet
	filename  string
	updatedAt time.Time

	data *Data
	// channel -> knob -> value
	knobs map[int]map[int]int

	showInst int
}

func New(filename string, m InstrumentSet) *Config {
	return &Config{
		m:        m,
		filename: filename,
		data:     new(Data),
		knobs:    make(map[int]map[int]int),
		showInst: -1, // all
	}
}

func (c *Config) InitMidiSynth() error {
	log.Printf("Synthesizer configuration: %v", c.filename)
	f, err := os.Open(c.filename)
	if err != nil {
		return fmt.Errorf("Error reading config file '%v': %w", c.filename, err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Cannot detect file modification time: %v", err)
	}
	modTime := fi.ModTime()
	if !modTime.After(c.updatedAt) {
		log.Printf("No need to update (%v <= %v)", modTime, c.updatedAt)
		return nil
	}
	c.updatedAt = modTime
	if err := yaml.NewDecoder(f).Decode(c.data); err != nil {
		return fmt.Errorf("Error parsing data: %w", err)
	}
	log.Printf("Loading configuration from %v", c.filename)
	if err := c.handleData(c.data, c.m); err != nil {
		return err
	}
	return nil
}

func (c *Config) StartUpdateLoop() error {
	callback := func() {
		log.Printf("Updating Midisynth...")
		if err := c.InitMidiSynth(); err != nil {
			log.Printf("Failed to update MidiSynth: %v", err)
		}
		log.Printf("Updating Midisynth... DONE.")
	}
	return watch.OnFileUpdate(c.filename, callback)
}

func (c *Config) handleData(data *Data, m InstrumentSet) error {
	var indexes []string
	for i := range data.Instruments {
		indexes = append(indexes, i)
	}
	sort.Strings(indexes)
	for _, inst := range indexes {
		instData := data.Instruments[inst]
		instIdx, err := strconv.Atoi(inst)
		if err != nil {
			return fmt.Errorf("Instrument index is not an integer: %v", inst)
		}
		c.log(instIdx, "Loading instrument %v...", instIdx)
		instr, err := c.handleInstrument(instIdx, &instData)
		if err != nil {
			return err
		}
		m.AddInstrument(instIdx, instr)
	}

	for name, sampleData := range data.Samples {
		c.log(-1, "Loading sampled instrument %v...", name)
		instr, err := c.handleSampleData(&sampleData)
		if err != nil {
			return fmt.Errorf("Failed to handle instrument %v: %v", name, err)
		}
		m.AddSampledInstrument(name, instr)
	}
	return nil
}

func (c *Config) handleInstrument(inst int, in *Instrument) (*instruments.Instrument, error) {
	wave, err := c.handleWave(inst, in.Wave)
	if err != nil {
		return nil, fmt.Errorf("Failed to handle wave: %w", err)
	}
	var fs []filters.Filter
	for _, f := range in.Filters {
		fx, err := c.handleFilter(inst, f)
		if err != nil {
			return nil, fmt.Errorf("Failed to handle filter: %w", err)
		}
		fs = append(fs, fx)
	}
	return instruments.NewInstrument(wave, fs...), nil
}

func (c *Config) handleSampleData(s *SampleData) (*instruments.Instrument, error) {
	sample, err := c.handleSample(s.Sample)
	if err != nil {
		return nil, fmt.Errorf("Failed to handle sample %v: %v", s.Sample, err)
	}
	var fs []filters.Filter
	for _, f := range s.Filters {
		fx, err := c.handleFilter(-1, f)
		if err != nil {
			return nil, fmt.Errorf("Failed to handle filter: %w", err)
		}
		fs = append(fs, fx)
	}
	return instruments.NewInstrument(sample, fs...), nil
}

func (c *Config) handleWave(inst int, wave string) (waves.Wave, error) {
	switch wave {
	case "sine":
		c.log(inst, "> Using Sine wave.")
		return waves.Sine, nil
	case "square":
		c.log(inst, "> Using Square wave.")
		return waves.Square, nil
	case "triangle":
		c.log(inst, "> Using Triangle wave.")
		return waves.Triangle, nil
	case "saw":
		c.log(inst, "> Using Sawtooth wave.")
		return waves.Saw, nil
	case "semisine":
		c.log(inst, "> Using Semisine wave.")
		return waves.SemiSine, nil
	}
	return nil, fmt.Errorf("Unknown wave: %v", wave)
}

func (c *Config) handleSample(sample string) (waves.Wave, error) {
	c.log(-1, "> Using sample '%v'", sample)
	data, err := static.Files.ReadFile(sample)
	if err != nil {
		return nil, err
	}
	return waves.ParseSample(data)
}

func (c *Config) handleFilter(instr int, f Filter) (filters.Filter, error) {
	if len(f) != 1 {
		return nil, fmt.Errorf("Filter description should contain exactly 1 element: %+v", f)
	}
	for name, opts := range f {
		switch name {
		case "adsr":
			return c.handleAdsr(instr, opts)
		case "delay":
			return c.handleDelay(instr, opts)
		case "distortion":
			return c.handleDistortion(instr, opts)
		case "vibrato":
			return c.handleVibrato(instr, opts)
		case "am":
			return c.handleAmplitudeModulation(instr, opts)
		case "timeshift":
			return c.handleTimeShift(instr, opts)
		case "harmonizer":
			return c.handleHarmonizer(instr, opts)
		case "flanger":
			return c.handleFlanger(instr, opts)
		}
		return nil, fmt.Errorf("Unknown filter: %v", name)
	}
	panic("unreachable")
}

func (c *Config) handleAdsr(instr int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(instr, "> Using ADSR filter...")
	var o []func(*filters.AdsrFilter)
	for param, value := range opts {
		c.log(instr, ">> with %v = %v", param, value)
		switch param {
		case "attackLevel":
			o = append(o, filters.AdsrAttackLevel(c.valueFloat64(instr, value)))
		case "decayLevel":
			o = append(o, filters.AdsrDecayLevel(c.valueFloat64(instr, value)))
		case "attackLen":
			o = append(o, filters.AdsrAttackLen(c.valueFloat64(instr, value)))
		case "decayLen":
			o = append(o, filters.AdsrDecayLen(c.valueFloat64(instr, value)))
		case "susteinLen":
			o = append(o, filters.AdsrSusteinLen(c.valueFloat64(instr, value)))
		case "releaseLen":
			o = append(o, filters.AdsrReleaseLen(c.valueFloat64(instr, value)))
		default:
			return nil, fmt.Errorf("Unknown ADSR parameter: %v", param)
		}
	}
	return filters.NewAdsrFilter(o...), nil
}

func (c *Config) handleDelay(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Delay filter")
	var o []func(*filters.DelayFilter)
	for param, value := range opts {
		switch param {
		case "interval":
			v := c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, v)
			o = append(o, filters.DelayInterval(v))
		case "times":
			v := value.(int)
			c.log(inst, ">> with %v = %v", param, v)
			o = append(o, filters.DelayTimes(v))
		case "fade":
			v := c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, v)
			o = append(o, filters.DelayFadeOut(v))
		default:
			return nil, fmt.Errorf("Unknown Delay parameter: %v", param)
		}
	}
	return filters.NewDelayFilter(o...), nil
}

func (c *Config) handleDistortion(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Distortion filter")
	value := 1.0
	for param, v := range opts {
		switch param {
		case "value":
			value = c.valueFloat64(inst, v)
			c.log(inst, ">> with %v = %v -> %v", param, v, value)
		default:
			return nil, fmt.Errorf("Unknown Distortion parameter: %v", param)
		}
	}
	return filters.NewDistortionFilter(value), nil
}

func (c *Config) handleVibrato(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Vibrato filter")
	var o []func(*filters.VibratoFilter)
	for param, value := range opts {
		switch param {
		case "wave":
			w, err := c.handleWave(inst, value.(string))
			if err != nil {
				return nil, err
			}
			c.log(inst, ">> with %v = %v", param, value)
			o = append(o, filters.VibratoCarrierWave(w))
		case "frequency":
			v := c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, v)
			o = append(o, filters.VibratoFrequency(v))
		case "amplitude":
			v := c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, v)
			o = append(o, filters.VibratoAmplitude(v))
		default:
			return nil, fmt.Errorf("Unknown Vibrato parameter: %v", param)
		}
	}
	return filters.NewVibrato(o...), nil
}

func (c *Config) handleAmplitudeModulation(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using AM (amplitude modulation) filter")
	var carrier waves.Wave = waves.Sine
	var freq float64
	amp := 1.0

	for param, value := range opts {
		switch param {
		case "wave":
			w, err := c.handleWave(inst, value.(string))
			if err != nil {
				return nil, err
			}
			carrier = w
			c.log(inst, ">> with %v = %v", param, value)
		case "frequency":
			freq = c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, freq)
		case "amplitude":
			amp = c.valueFloat64(inst, value)
			c.log(inst, ">> with %v = %v -> %v", param, value, amp)
		default:
			return nil, fmt.Errorf("Unknown AM parameter: %v", param)
		}

	}

	return filters.NewRing(carrier, freq, amp), nil
}

func (c *Config) handleTimeShift(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Time Shift filter")
	var o []func(*filters.TimeShift)
	for param, value := range opts {
		c.log(inst, ">> with %v = %v", param, value)
		switch param {
		case "wave":
			w, err := c.handleWave(inst, value.(string))
			if err != nil {
				return nil, err
			}
			o = append(o, filters.TimeShiftCarrierWave(w))
		case "frequency":
			o = append(o, filters.TimeShiftFrequency(c.valueFloat64(inst, value)))
		case "amplitude":
			o = append(o, filters.TimeShiftAmplitude(c.valueFloat64(inst, value)))
		default:
			return nil, fmt.Errorf("Unknown Time Shift parameter: %v", param)
		}
	}
	return filters.NewTimeShift(o...), nil
}

func (c *Config) handleHarmonizer(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Harmonizer filter")
	keys := make([]string, 0, len(opts))
	for param := range opts {
		keys = append(keys, param)
	}
	sort.Strings(keys)
	var o []func(*filters.Harmonizer)
	for _, param := range keys {
		n, err := strconv.Atoi(param)
		if err != nil {
			return nil, fmt.Errorf("Incorrect Harmonizer param: %v", param)
		}
		value := opts[param]
		v := c.valueFloat64(inst, value)
		c.log(inst, "  >> with %v = %v (%v)", n, v, value)
		o = append(o, filters.Harmonic(n, v))
	}
	return filters.NewHarmonizer(o...), nil
}

func (c *Config) handleFlanger(inst int, opts map[string]interface{}) (filters.Filter, error) {
	c.log(inst, "> Using Flanger filter")
	var o []func(*filters.Flanger)
	for param, value := range opts {
		switch param {
		case "frequency":
			v := c.valueFloat64(inst, value)
			c.log(inst, "  >> with %v = %v", param, v)
			o = append(o, filters.FlangerFreq(v))
		default:
			return nil, fmt.Errorf("Unknown Flanger parameter: %v", param)
		}
	}

	return filters.NewFlanger(o...), nil
}

func (c *Config) log(inst int, format string, args ...interface{}) {
	if c.showInst < 0 || c.showInst == inst {
		log.Printf(format, args...)
	}
}

func (c *Config) valueFloat64(instr int, x interface{}) float64 {
	switch a := x.(type) {
	case float64:
		return a
	case int:
		return float64(a)
	case map[string]interface{}:
		// parameter:
		//   knob: 1
		//   default: 10.0
		//   delta: 0.1
		knob, ok := a["knob"].(int)
		if !ok {
			panic(fmt.Errorf("Integer value 'knob' not found in %v", a))
		}
		def, ok := a["default"].(float64)
		if !ok {
			panic(fmt.Errorf("Float value 'default' not found in %v", a))
		}
		delta, ok := a["delta"].(float64)
		if !ok {
			panic(fmt.Errorf("Float value 'delta' not found in %v", a))
		}
		return c.knobValue(instr, knob, def, delta)
	default:
		panic(fmt.Errorf("Not an integer value: %v", x))
	}
}

func (c *Config) knobValue(inst int, knob int, def float64, delta float64) float64 {
	ik, ok := c.knobs[inst]
	if !ok {
		return def
	}
	kv, ok := ik[knob]
	if !ok {
		return def
	}
	return def + (float64(kv) * delta)
}

func (c *Config) Up(inst, knob int) {
	ik, ok := c.knobs[inst]
	if !ok {
		c.knobs[inst] = map[int]int{
			knob: 1,
		}
		return
	}
	ik[knob] += 1
	c.knobs[inst] = ik
	log.Printf("Up: knobs = %+v", c.knobs)
	c.showInst = inst
	if err := c.handleData(c.data, c.m); err != nil {
		log.Printf("Cannot update configuration: %v", err)
	}
	c.showInst = -1
}

func (c *Config) Down(inst, knob int) {
	ik, ok := c.knobs[inst]
	if !ok {
		c.knobs[inst] = map[int]int{
			knob: -1,
		}
		return
	}
	ik[knob] -= 1
	c.knobs[inst] = ik
	log.Printf("Down: knobs = %+v", c.knobs)
	c.showInst = inst
	if err := c.handleData(c.data, c.m); err != nil {
		log.Printf("Cannot update configuration: %v", err)
	}
	c.showInst = -1
}
