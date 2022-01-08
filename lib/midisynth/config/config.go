package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Config struct {
}

func (c *Config) InitMidiSynth(filename string, m *midisynth.MidiSynth) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading config file '%v': %w", filename, err)
	}
	var data Data
	if err := json.Unmarshal(f, &data); err != nil {
		return fmt.Errorf("Error parsing data: %w\n%s", err, f)
	}
	log.Printf("Loading configuration from %v", filename)
	if err := c.handleData(&data, m); err != nil {
		return err
	}
	return nil
}

func (c *Config) handleData(data *Data, m *midisynth.MidiSynth) error {
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
			o = append(o, filters.AdsrAttackLevel(value.(float64)))
		case "decayLevel":
			o = append(o, filters.AdsrDecayLevel(value.(float64)))
		case "attackLen":
			o = append(o, filters.AdsrAttackLen(value.(float64)))
		case "decayLen":
			o = append(o, filters.AdsrDecayLen(value.(float64)))
		case "susteinLen":
			o = append(o, filters.AdsrSusteinLen(value.(float64)))
		case "releaseLen":
			o = append(o, filters.AdsrReleaseLen(value.(float64)))
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
			o = append(o, filters.DelayInterval(value.(float64)))
		case "times":
			o = append(o, filters.DelayTimes(int(value.(float64))))
		case "fade":
			o = append(o, filters.DelayFadeOut(value.(float64)))
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
			value = v.(float64)
		default:
			return nil, fmt.Errorf("Unknown Distortion parameter: %v", param)
		}
	}
	return filters.NewDistortionFilter(value), nil
}
