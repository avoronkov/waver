package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Synth interface {
	PlayNote(instr int, octave int, note string, durationBits int, amp float64) error
}

type Proc struct {
	synth Synth
}

// []string{"24:0", "Note", "off", "1,", "note", "67,", "velocity", "64"}
// []string{"24:0", "Control", "change", "1,", "controller", "2,", "value", "12"}
func (p *Proc) handleLine(line string) error {
	fields := strings.Fields(line)
	log.Printf("> %#v", fields)
	switch fmt.Sprintf("%v %v", fields[1], fields[2]) {
	case "Note on":
		log.Printf("Note on")
		err := p.handleNoteOn(fields)
		if err != nil {
			log.Printf("p.handleNoteOn failed: %v", err)
		}
	case "Note off":
		log.Printf("Note off")

	}
	return nil
}

func (p *Proc) channel(fields []string) (int, error) {
	return strconv.Atoi(strings.TrimRight(fields[3], ","))
}

type Key struct {
	channel  int
	note     int
	velocity int
}

func (p *Proc) handleNoteOn(fields []string) error {
	key, err := p.parseNote(fields)
	if err != nil {
		return err
	}
	on, ok := KeyMap[key.note]
	if !ok {
		panic(fmt.Errorf("Key not found in table: %v", key.note))
	}
	return p.synth.PlayNote(key.channel, on.Octave, on.Note, 4, 0.2)
}

func (p *Proc) parseNote(fields []string) (key *Key, err error) {
	key = &Key{}
	key.channel, err = strconv.Atoi(strings.TrimRight(fields[3], ","))
	if err != nil {
		return nil, err
	}
	if fields[4] == "note" {
		key.note, err = strconv.Atoi(strings.TrimRight(fields[5], ","))
	} else {
		return nil, fmt.Errorf("Cannot parse 'note' section")
	}
	if fields[6] == "velocity" {
		key.velocity, err = strconv.Atoi(strings.TrimRight(fields[7], ","))
	} else {
		return nil, fmt.Errorf("Cannot parse 'velocity' secion")
	}

	return
}
