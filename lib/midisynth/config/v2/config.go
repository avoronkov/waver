package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
}

func (c *Config) InitMidiSynth(filename string, m *midisynth.MidiSynth) error {
	log.Printf("Synthesizer configuration: %v", filename)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading config file '%v': %w", filename, err)
	}
	var data Data
	if err := yaml.Unmarshal(f, &data); err != nil {
		return fmt.Errorf("Error parsing data: %w\n%s", err, f)
	}
	log.Printf("Loading configuration from %v", filename)
	if err := c.handleData(&data, m); err != nil {
		return err
	}
	return nil
}

func (c *Config) handleData(data *Data, m *midisynth.MidiSynth) error {
	log.Printf("Data: %+v", data)
	for inst, instData := range data.Instruments {
		instIdx, err := strconv.Atoi(inst)
		if err != nil {
			return fmt.Errorf("Instrument index is not an integer: %v", inst)
		}
		log.Printf("Loading instrument %v...", instIdx)
		instr, err := c.handleInstrument(&instData)
		if err != nil {
			return err
		}
		m.AddInstrument(instIdx, instr)
	}

	for name, sampleData := range data.Samples {
		log.Printf("Loading sampled instrument %v...", name)
		instr, err := c.handleSampleData(&sampleData)
		if err != nil {
			return fmt.Errorf("Failed to handle instrument %v: %v", name, err)
		}
		m.AddSampledInstrument(name, instr)
	}
	return nil
}

func (c *Config) handleInstrument(in *Instrument) (*instruments.Instrument, error) {
	wave, err := c.handleWave(in.Wave)
	if err != nil {
		return nil, fmt.Errorf("Failed to handle wave: %w", err)
	}
	var fs []filters.Filter
	for _, f := range in.Filters {
		fx, err := c.handleFilter(f)
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
		fx, err := c.handleFilter(f)
		if err != nil {
			return nil, fmt.Errorf("Failed to handle filter: %w", err)
		}
		fs = append(fs, fx)
	}
	return instruments.NewInstrument(sample, fs...), nil
}

func (c *Config) handleWave(wave string) (waves.Wave, error) {
	switch wave {
	case "sine":
		log.Printf("> Using Sine wave.")
		return &waves.Sine{}, nil
	case "square":
		log.Printf("> Using Square wave.")
		return &waves.Square{}, nil
	case "triangle":
		log.Printf("> Using Triangle wave.")
		return &waves.Triangle{}, nil
	case "saw":
		log.Printf("> Using Sawtooth wave.")
		return &waves.Saw{}, nil
	}
	return nil, fmt.Errorf("Unknown wave: %v", wave)
}

func (c *Config) handleSample(sample string) (waves.Wave, error) {
	log.Printf("> Using sample '%v'", sample)
	return waves.ReadSample(sample)
}

func (c *Config) handleFilter(f Filter) (filters.Filter, error) {
	if len(f) != 1 {
		return nil, fmt.Errorf("Filter description should contain exactly 1 element: %+v", f)
	}
	for name, opts := range f {
		switch name {
		case "adsr":
			return c.handleAdsr(opts)
		case "delay":
			return c.handleDelay(opts)
		case "distortion":
			return c.handleDistortion(opts)
		case "vibrato":
			return c.handleVibrato(opts)
		case "am":
			return c.handleAmplitudeModulation(opts)
		case "timeshift":
			return c.handleTimeShift(opts)
		}
		return nil, fmt.Errorf("Unknown filter: %v", name)
	}
	panic("unreachable")
}

func (c *Config) handleAdsr(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using ADSR filter...")
	var o []func(*filters.AdsrFilter)
	for param, value := range opts {
		log.Printf(">> with %v = %v", param, value)
		switch param {
		case "attackLevel":
			o = append(o, filters.AdsrAttackLevel(valueFloat64(value)))
		case "decayLevel":
			o = append(o, filters.AdsrDecayLevel(valueFloat64(value)))
		case "attackLen":
			o = append(o, filters.AdsrAttackLen(valueFloat64(value)))
		case "decayLen":
			o = append(o, filters.AdsrDecayLen(valueFloat64(value)))
		case "susteinLen":
			o = append(o, filters.AdsrSusteinLen(valueFloat64(value)))
		case "releaseLen":
			o = append(o, filters.AdsrReleaseLen(valueFloat64(value)))
		default:
			return nil, fmt.Errorf("Unknown ADSR parameter: %v", param)
		}
	}
	return filters.NewAdsrFilter(o...), nil
}

func (c *Config) handleDelay(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using Delay filter")
	var o []func(*filters.DelayFilter)
	for param, value := range opts {
		log.Printf(">> with %v = %v", param, value)
		switch param {
		case "interval":
			o = append(o, filters.DelayInterval(valueFloat64(value)))
		case "times":
			o = append(o, filters.DelayTimes(value.(int)))
		case "fade":
			o = append(o, filters.DelayFadeOut(valueFloat64(value)))
		default:
			return nil, fmt.Errorf("Unknown Delay parameter: %v", param)
		}
	}
	return filters.NewDelayFilter(o...), nil
}

func (c *Config) handleDistortion(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using Distortion filter")
	value := 1.0
	for param, v := range opts {
		log.Printf(">> with %v = %v", param, v)
		switch param {
		case "value":
			value = valueFloat64(v)
		default:
			return nil, fmt.Errorf("Unknown Distortion parameter: %v", param)
		}
	}
	return filters.NewDistortionFilter(value), nil
}

func (c *Config) handleVibrato(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using Vibrato filter")
	var o []func(*filters.VibratoFilter)
	for param, value := range opts {
		log.Printf(">> with %v = %v", param, value)
		switch param {
		case "wave":
			w, err := c.handleWave(value.(string))
			if err != nil {
				return nil, err
			}
			o = append(o, filters.VibratoCarrierWave(w))
		case "frequency":
			o = append(o, filters.VibratoFrequency(valueFloat64(value)))
		case "amplitude":
			o = append(o, filters.VibratoAmplitude(valueFloat64(value)))
		default:
			return nil, fmt.Errorf("Unknown Vibrato parameter: %v", param)
		}
	}
	return filters.NewVibrato(o...), nil
}

func (c *Config) handleAmplitudeModulation(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using AM (amplitude modulation) filter")
	var carrier waves.Wave = &waves.Sine{}
	var freq float64
	amp := 1.0

	for param, value := range opts {
		log.Printf(">> with %v = %v", param, value)
		switch param {
		case "wave":
			w, err := c.handleWave(value.(string))
			if err != nil {
				return nil, err
			}
			carrier = w
		case "frequency":
			freq = valueFloat64(value)
		case "amplitude":
			amp = valueFloat64(value)
		default:
			return nil, fmt.Errorf("Unknown AM parameter: %v", param)
		}

	}

	return filters.NewRing(carrier, freq, amp), nil
}

func (c *Config) handleTimeShift(opts map[string]interface{}) (filters.Filter, error) {
	log.Printf("> Using Time Shift filter")
	var o []func(*filters.TimeShift)
	for param, value := range opts {
		log.Printf(">> with %v = %v", param, value)
		switch param {
		case "wave":
			w, err := c.handleWave(value.(string))
			if err != nil {
				return nil, err
			}
			o = append(o, filters.TimeShiftCarrierWave(w))
		case "frequency":
			o = append(o, filters.TimeShiftFrequency(valueFloat64(value)))
		case "amplitude":
			o = append(o, filters.TimeShiftAmplitude(valueFloat64(value)))
		default:
			return nil, fmt.Errorf("Unknown Time Shift parameter: %v", param)
		}
	}
	return filters.NewTimeShift(o...), nil
}

func valueFloat64(x interface{}) float64 {
	switch a := x.(type) {
	case float64:
		return a
	case int:
		return float64(a)
	default:
		panic(fmt.Errorf("Not an integer value: %v", x))
	}
}
