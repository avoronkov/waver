package udp

import (
	"errors"
	"fmt"
	"log"
	"net"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
)

type UdpInput struct {
	port int

	listener net.PacketConn
}

var _ signals.Input = (*UdpInput)(nil)

func New(port int) *UdpInput {
	return &UdpInput{
		port: port,
	}
}

func (u *UdpInput) Start(ch chan<- *signals.Signal) (err error) {
	log.Printf("Starting UDP listener on port %v...", u.port)
	u.listener, err = net.ListenPacket("udp", fmt.Sprintf(":%v", u.port))
	if err != nil {
		return fmt.Errorf("Starting UDP server failed: %w", err)
	}
	log.Printf("Listening to UDP on localhost:%v", u.port)
	go func(pc net.PacketConn) {
	L:
		for {
			buff := make([]byte, 64)
			n, _, err := pc.ReadFrom(buff)
			if errors.Is(err, net.ErrClosed) {
				log.Printf("[ERROR] ErrClosed")
				break L
			}
			if err != nil {
				log.Printf("[ERROR] %v (%T, %v)", err, err, err.(*net.OpError).Unwrap())
				continue
			}
			sig, err := parseMessage(buff[:n])
			if err != nil {
				log.Printf("[ERROR] %v", err)
				continue
			}
			if sig != nil {
				ch <- sig
			}
		}
	}(u.listener)
	return nil
}

func (u *UdpInput) Close() error {
	return u.listener.Close()
}

func parseMessage(msg []byte) (*signals.Signal, error) {
	if len(msg) < 3 {
		return nil, nil
	}
	inst := parseValue(msg[0])
	octave := int(msg[1] - '0')
	note := string(msg[2])
	amp := 0.5
	if len(msg) >= 4 {
		amp = 0.1 * float64(parseValue(msg[3]))
	}
	dur := 4
	if len(msg) >= 5 {
		// Evaluate duration in bits (1/4 tempo)
		dur = parseDuration(msg[4:])
	}

	if inst == 35 { // 'z'
		// handle samples
		return &signals.Signal{
			Sample:       string(msg[1:3]),
			DurationBits: dur,
			Amp:          amp,
		}, nil
	}

	// regular notes
	return &signals.Signal{
		Instrument:   inst,
		Octave:       octave,
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
