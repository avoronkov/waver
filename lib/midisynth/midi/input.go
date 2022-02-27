package midi

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
)

type Input struct {
	midiPort    int
	dumpProcess *exec.Cmd

	keyMap map[int]OctaveNote
}

var _ signals.Input = (*Input)(nil)

func NewInput(midiPort int) *Input {
	return &Input{
		midiPort: midiPort,
		// TODO configure keymap
		keyMap: KeyMap,
	}
}

func (i *Input) Start(ch chan<- *signals.Signal) (err error) {
	log.Printf("Starting aseqdump process (-p %v)", i.midiPort)
	dumpProcess, reader, err := aseqdump(i.midiPort)
	if err != nil {
		return err
	}
	i.dumpProcess = dumpProcess

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	go func() {
		log.Printf("Scanning aseqdump outpupt")
		for scanner.Scan() {
			text := scanner.Text()
			log.Printf("[MIDI] Got event: %v", text)
			sig, err := parseLine(i.keyMap, text)
			if err != nil {
				log.Printf("[ERROR] parseLine failed: %v", text)
				continue
			}
			if sig != nil {
				ch <- sig
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}()
	return nil
}

func (i *Input) Close() error {
	if i.dumpProcess != nil {
		_ = i.dumpProcess.Process.Signal(syscall.SIGINT)
	}
	return nil
}

func parseLine(keyMap map[int]OctaveNote, line string) (*signals.Signal, error) {
	fields := strings.Fields(line)
	log.Printf("> %#v", fields)
	if len(fields) < 3 {
		return nil, nil
	}
	switch fmt.Sprintf("%v %v", fields[1], fields[2]) {
	case "Note on":
		log.Printf("Note on")
		key, err := parseNote(keyMap, fields)
		if err != nil {
			log.Printf("parseNote failed: %v (%v)", err, fields)
		}
		return key, nil
	case "Note off":
		log.Printf("Note off")
		key, err := parseNote(keyMap, fields)
		if err != nil {
			log.Printf("parseNote failed: %v (%v)", err, fields)
		}
		key.Stop = true
		return key, nil
	case "Control change":
		log.Printf("Control change (NIY)")
		/*
			err := p.handleControlChange(fields)
			if err != nil {
				log.Printf("p.handleControlChange failed: %v", err)
			}
		*/
		return nil, nil
	}

	return nil, nil
}

func parseNote(keyMap map[int]OctaveNote, fields []string) (key *signals.Signal, err error) {
	key = &signals.Signal{
		Manual: true,
	}
	key.Instrument, err = strconv.Atoi(strings.TrimRight(fields[3], ","))
	if err != nil {
		return nil, err
	}
	key.Instrument++
	if fields[4] == "note" {
		noteIdx, err := strconv.Atoi(strings.TrimRight(fields[5], ","))
		if err != nil {
			return nil, err
		}
		on, ok := keyMap[noteIdx]
		if !ok {
			return nil, fmt.Errorf("Unknown note index: %v", noteIdx)
		}
		key.Octave = on.Octave
		key.Note = on.Note
	} else {
		return nil, fmt.Errorf("Cannot parse 'note' section")
	}

	// TODO parse velocity
	key.Amp = 0.1

	return key, nil
}
