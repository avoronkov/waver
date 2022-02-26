package seq

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type Sequencer struct {
	tempo int
	port  int

	conn *net.UDPConn
	sigs chan os.Signal

	started bool
	current []types.Signaler
	next    []types.Signaler
}

const defaultPort = 49161

func NewSequencer() *Sequencer {
	s := &Sequencer{
		tempo: 120,
		port:  defaultPort,
		sigs:  make(chan os.Signal, 1),
	}
	return s
}

func (s *Sequencer) Run(funcs ...types.Signaler) error {
	if err := s.init(); err != nil {
		return err
	}
	s.current = funcs
	if err := s.run(); err != nil {
		return err
	}
	return s.close()
}

func (s *Sequencer) Add(sig types.Signaler) {
	s.next = append(s.next, sig)
}

func (s *Sequencer) Commit() error {
	s.current = s.next
	s.next = nil
	if !s.started {
		s.started = true
		if err := s.init(); err != nil {
			return err
		}
		go s.run()
	}
	return nil
}

func (s *Sequencer) init() error {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%v", s.port))
	if err != nil {
		return err
	}
	s.conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	signal.Notify(s.sigs, syscall.SIGINT, syscall.SIGTERM)

	return nil
}

func (s *Sequencer) run() error {
	delay := time.Duration((15.0 / float64(s.tempo)) * float64(time.Second))
	currentDelay := 0 * time.Millisecond
	var bit int64
	for {
		select {
		case <-time.After(currentDelay):
			start := time.Now()
			s.processFuncs(bit, s.current)
			dt := time.Since(start)
			currentDelay = delay - dt
		case sig := <-s.sigs:
			log.Printf("Got signal %v. Terminating.", sig)
			return s.close()
		}
		bit++
	}
}

func (s *Sequencer) processFuncs(bit int64, funcs []types.Signaler) {
	for _, fn := range funcs {
		signals := fn.Eval(bit, types.Context{})
		for _, sig := range signals {
			log.Print(sig)
			fmt.Fprintf(s.conn, sig)
		}
	}
}

func (s *Sequencer) close() error {
	return s.conn.Close()
}

var DefaultSequencer = NewSequencer()

func Run(funcs ...types.Signaler) error {
	return DefaultSequencer.Run(funcs...)
}
