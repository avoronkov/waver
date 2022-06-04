//go:build !js
// +build !js

package config

import (
	"log"

	"github.com/avoronkov/waver/lib/watch"
)

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
