package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/static"

	yaml "gopkg.in/yaml.v3"
)

type InstrumentSet interface {
	AddInstrument(n string, in *instruments.Instrument)
	// AddSampledInstrument(name string, in *instruments.Instrument)
}

type Config struct {
	m         InstrumentSet
	filename  string
	updatedAt time.Time

	data *Data
	// channel -> knob -> value
	knobs map[string]map[int]int

	showInst string
}

func New(filename string, m InstrumentSet) *Config {
	return &Config{
		m:        m,
		filename: filename,
		data:     new(Data),
		knobs:    make(map[string]map[int]int),
		showInst: "", // all
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
	log.Printf("Loading configuration from %v", c.filename)
	return c.updateReader(f)
}

func (c *Config) UpdateData(data []byte) error {
	return c.updateReader(bytes.NewReader(data))
}

func (c *Config) updateReader(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(c.data); err != nil {
		return fmt.Errorf("Error parsing data: %w", err)
	}
	if err := c.handleData(c.data, c.m); err != nil {
		return err
	}
	return nil
}

func (c *Config) handleData(data *Data, m InstrumentSet) error {
	var indexes []string
	for i := range data.Instruments {
		indexes = append(indexes, i)
	}
	sort.Strings(indexes)
	for _, inst := range indexes {
		instData := data.Instruments[inst]
		c.log(inst, "Loading instrument %v...", inst)
		instr, err := ParseInstrument(instData.Wave, instData.Filters)
		if err != nil {
			return err
		}
		m.AddInstrument(inst, instr)
	}

	for name, sampleData := range data.Samples {
		c.log(name, "Loading sampled instrument %v...", name)
		instr, err := ParseSample(sampleData.Sample, sampleData.Filters)
		if err != nil {
			return fmt.Errorf("Failed to handle instrument %v: %v", name, err)
		}
		m.AddInstrument(name, instr)
	}
	return nil
}

func ParseSample(file string, filtersData []map[string]any) (*instruments.Instrument, error) {
	sample, err := handleSample(file)
	if err != nil {
		return nil, err
	}
	var fs []filters.Filter
	for _, f := range filtersData {
		fx, err := handleFilter(f)
		if err != nil {
			return nil, fmt.Errorf("Failed to handle filter: %w", err)
		}
		fs = append(fs, fx)
	}
	return instruments.NewInstrument(sample, fs...), nil
}

func handleSample(sample string) (waves.Wave, error) {
	data, err := static.Files.ReadFile(sample)
	if err != nil {
		return nil, err
	}
	return waves.ParseSample(data)
}

func ParseInstrument(waveName string, filtersData []map[string]any) (*instruments.Instrument, error) {
	w, ok := waves.Waves[waveName]
	if !ok {
		return nil, fmt.Errorf("Unknown wave: %v", waveName)
	}
	var fs []filters.Filter
	for _, f := range filtersData {
		fx, err := handleFilter(f)
		if err != nil {
			return nil, fmt.Errorf("Failed to handle filter: %w", err)
		}
		fs = append(fs, fx)
	}
	return instruments.NewInstrument(w, fs...), nil
}

func handleFilter(f map[string]any) (filters.Filter, error) {
	if len(f) != 1 {
		return nil, fmt.Errorf("Filter description should contain exactly 1 element: %+v", f)
	}
	for name, opts := range f {
		fc, ok := filters.Filters[name]
		if !ok {
			return nil, fmt.Errorf("Unknown filter: %v", name)
		}
		return fc.Create(opts)
	}
	panic("unreachable")
}

func (c *Config) log(inst string, format string, args ...interface{}) {
	if c.showInst == "" || c.showInst == inst {
		log.Printf(format, args...)
	}
}

func (c *Config) knobValue(inst string, knob int, def float64, delta float64) float64 {
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

func (c *Config) Up(inst string, knob int) {
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
	c.showInst = ""
}

func (c *Config) Down(inst string, knob int) {
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
	c.showInst = ""
}
