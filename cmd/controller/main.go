package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func main() {
	flag.Parse()

	log.Printf("aconnect -i")
	check(aconnectI())
	log.Printf("aseqdump -p %v", p)
	dumpProcess, reader, err := aseqdump(p)
	check(err)

	scale := notes.NewStandard()
	m, err := midisynth.NewMidiSynth(wav.Default, scale, 49161)
	check(err)

	cfg := &config.Config{}
	check(cfg.InitMidiSynth(configPath, m))

	// Experimental section

	m.AddInstrument(9, instruments.NewInstrument(
		&waves.Sine{},
		filters.NewVibrato(&waves.Sine{}, 10.0, 0.05),
		filters.NewAdsrFilter(),
	))

	m.AddInstrument(8, instruments.NewInstrument(
		&waves.Sine{},
		filters.NewTimeShift(10.0, 0.01),
		filters.NewAdsrFilter(),
	))

	m.AddInstrument(7, instruments.NewInstrument(
		&waves.Triangle{},
		filters.NewRing(&waves.Sine{}, 4.0),
		filters.NewAdsrFilter(),
	))

	// .
	proc := NewProc(m)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	midi := make(chan string)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			midi <- text
		}
		check(scanner.Err())
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

L:
	for {
		select {
		case text := <-midi:
			go func() {
				if err := proc.handleLine(text); err != nil {
					log.Printf("[WARN] %v", err)
				}
			}()
		case <-sigs:
			log.Printf("Interupting...")
			dumpProcess.Process.Signal(syscall.SIGINT)
			break L
		}
	}

	log.Printf("Done")
}

func aconnectI() error {
	cmd := exec.Command("aconnect", "-i")
	return cmd.Run()
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
