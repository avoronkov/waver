package midi

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Synth interface {
	PlayNote(instr int, octave int, note string, durationBits int, amp float64) error
	PlayNoteControlled(instr int, octave int, note string, amp float64) (func(), error)
}

type Proc struct {
	synth    Synth
	ch       chan<- string
	midiPort int

	notesReleases map[int]func()
	dumpProcess   *exec.Cmd

	keyMap map[int]OctaveNote

	dispatcher Dispatcher
}

func NewProc(synth Synth, midiPort int, opts ...func(*Proc)) *Proc {
	p := &Proc{
		synth:         synth,
		midiPort:      midiPort,
		notesReleases: make(map[int]func()),
		keyMap:        KeyMap,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Proc) Start() error {
	log.Printf("Starting aseqdump process (-p %v)", p.midiPort)
	dumpProcess, reader, err := aseqdump(p.midiPort)
	if err != nil {
		return err
	}
	p.dumpProcess = dumpProcess

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	go func() {
		log.Printf("Scanning aseqdump outpupt")
		for scanner.Scan() {
			text := scanner.Text()
			log.Printf("[MIDI] Got event: %v", text)
			p.ch <- text
		}
		if err := scanner.Err(); err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}()
	return nil
}

func (p *Proc) Close() error {
	if p.dumpProcess != nil {
		p.dumpProcess.Process.Signal(syscall.SIGINT)
	}
	return nil
}

// []string{"24:0", "Note", "off", "1,", "note", "67,", "velocity", "64"}
// []string{"24:0", "Control", "change", "1,", "controller", "2,", "value", "12"}
func (p *Proc) HandleLine(line string) error {
	fields := strings.Fields(line)
	log.Printf("> %#v", fields)
	if len(fields) < 3 {
		return nil
	}
	switch fmt.Sprintf("%v %v", fields[1], fields[2]) {
	case "Note on":
		log.Printf("Note on")
		err := p.handleNoteOn(fields)
		if err != nil {
			log.Printf("p.handleNoteOn failed: %v", err)
		}
	case "Note off":
		log.Printf("Note off")
		err := p.handleNoteOff(fields)
		if err != nil {
			log.Printf("p.handleNoteOff failed: %v", err)
		}
	case "Control change":
		log.Printf("Control change")
		err := p.handleControlChange(fields)
		if err != nil {
			log.Printf("p.handleControlChange failed: %v", err)
		}
	}
	return nil
}

type Key struct {
	channel  int
	note     int
	velocity int
}

// 24:0   Note on                 0, note 60, velocity 100
func (p *Proc) handleNoteOn(fields []string) error {
	key, err := p.parseNote(fields)
	if err != nil {
		return err
	}
	on, ok := p.keyMap[key.note]
	if !ok {
		panic(fmt.Errorf("Key not found in table: %v", key.note))
	}
	stop, err := p.synth.PlayNoteControlled(key.channel+1, on.Octave, on.Note, 0.2)
	if err != nil {
		return err
	}
	p.notesReleases[key.note] = stop
	return nil
}

// 24:0   Note off                0, note 62, velocity 64
func (p *Proc) handleNoteOff(fields []string) error {
	key, err := p.parseNote(fields)
	if err != nil {
		return err
	}
	if stop, ok := p.notesReleases[key.note]; ok && stop != nil {
		stop()
		p.notesReleases[key.note] = nil
	}
	return nil
}

// 24:0   Control change          0, controller 51, value 0
func (p *Proc) handleControlChange(fields []string) error {
	if p.dispatcher == nil {
		return nil
	}

	channel, err := strconv.Atoi(strings.TrimRight(fields[3], ","))
	if err != nil {
		return err
	}
	var controller, value int
	if fields[4] == "controller" {
		controller, err = strconv.Atoi(strings.TrimRight(fields[5], ","))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Cannot find controller section")
	}

	if fields[6] == "value" {
		value, err = strconv.Atoi(strings.TrimRight(fields[7], ","))
		if err != nil {
			return err
		}
		if value != 1 && value != 127 {
			return fmt.Errorf("Knob %v is not in relative mode: %v", controller, value)
		}
	} else {
		return fmt.Errorf("Cannot find 'value' section")
	}
	if value == 1 {
		p.dispatcher.Up(channel+1, controller)
	} else if value == 127 {
		p.dispatcher.Down(channel+1, controller)
	}
	return nil
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

func aseqdump(p int) (*exec.Cmd, io.Reader, error) {
	cmd := exec.Command("aseqdump", "-p", strconv.Itoa(p))
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, reader, nil
}
