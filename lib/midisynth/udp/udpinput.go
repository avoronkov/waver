package udp

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/notes"
)

type UdpInput struct {
	port int

	listener net.PacketConn

	scale notes.Scale
}

var _ signals.Input = (*UdpInput)(nil)

func New(port int, scale notes.Scale) *UdpInput {
	return &UdpInput{
		port:  port,
		scale: scale,
	}
}

func (u *UdpInput) Run(ch chan<- *signals.Signal) (err error) {
	log.Printf("Starting UDP listener on port %v...", u.port)
	u.listener, err = net.ListenPacket("udp", fmt.Sprintf(":%v", u.port))
	if err != nil {
		return fmt.Errorf("Starting UDP server failed: %w", err)
	}
	log.Printf("Listening to UDP on localhost:%v", u.port)
	for {
		buff := make([]byte, 64)
		n, _, err := u.listener.ReadFrom(buff)
		if errors.Is(err, net.ErrClosed) {
			return fmt.Errorf("UDP socked unexpectedly closed: %w", err)
		}
		if err != nil {
			log.Printf("[ERROR] %v (%T, %v)", err, err, err.(*net.OpError).Unwrap())
			continue
		}
		sig, err := ParseMessage(buff[:n], u.scale)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			continue
		}
		if sig != nil {
			ch <- sig
		}
	}
}

func (u *UdpInput) Close() error {
	return u.listener.Close()
}

func ParseMessage(msg []byte, scale notes.Scale) (*signals.Signal, error) {
	if len(msg) < 3 {
		return nil, nil
	}
	inst := msg[0]
	octave := int(msg[1] - '0')
	nt := string(msg[2])

	amp := 0.5
	if len(msg) >= 4 {
		amp = 0.1 * float64(parseValue(msg[3]))
	}
	dur := 4
	if len(msg) >= 5 {
		// Evaluate duration in bits (1/4 tempo)
		dur = parseDuration(msg[4:])
	}

	if inst == 'z' { // 'z'
		// handle samples
		return &signals.Signal{
			Sample:       string(msg[1:3]),
			DurationBits: dur,
			Amp:          amp,
		}, nil
	}

	note, ok := scale.Note(octave, nt)
	if !ok {
		return nil, fmt.Errorf("Cannot parse UDP note: %v%v", octave, nt)
	}

	// regular notes
	return &signals.Signal{
		Instrument:   string(inst),
		Note:         note,
		DurationBits: dur,
		Amp:          amp,
	}, nil
}

func parseValue(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'a' && b <= 'z' {
		return 10 + int(b-'a')
	}
	if b >= 'A' && b <= 'Z' {
		return 10 + int(b-'A')
	}
	return 0
}

func parseDuration(b []byte) int {
	if len(b) < 1 {
		panic("Empty duration")
	}
	v := parseValue(b[0])
	if len(b) >= 2 {
		v *= parseValue(b[1])
	}
	return v
}
