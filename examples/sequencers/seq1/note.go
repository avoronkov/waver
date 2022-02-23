package main

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq"
)

func newInt(n int) *int {
	return &n
}

func newRune(r rune) *rune {
	return &r
}

type note struct {
	instr  *int
	octave *int
	note   *rune
	amp    *int
	dur    *int
}

func Note(opts ...func(*note)) seq.Signaler {
	n := &note{
		amp: newInt(5),
		dur: newInt(4),
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

var _ seq.Signaler = (*note)(nil)

func (n *note) Eval(bit int64, ctx seq.Context) []string {
	instr := n.instr
	if instr == nil {
		in, ok := ctx["instr"].(int)
		if !ok {
			return nil
		}
		instr = &in
	}
	// octave
	octave := n.octave
	if octave == nil {
		oc, ok := ctx["octave"].(int)
		if !ok {
			return nil
		}
		octave = &oc
	}
	// note
	nt := n.note
	if nt == nil {
		n, ok := ctx["note"].(rune)
		if !ok {
			return nil
		}
		nt = &n
	}
	amp := n.amp
	if amp == nil {
		if n, ok := ctx["amp"].(int); ok {
			amp = &n
		}
	}

	result := fmt.Sprintf("%c%c%c", toRune(*instr), toRune(*octave), *nt)
	if amp != nil {
		result += string(toRune(*amp))
	}

	// TODO amp dur
	return []string{result}
}

func toRune(n int) rune {
	if n >= 0 && n < 10 {
		return '0' + rune(n)
	}
	if n >= 10 {
		return 'A' + rune(n) - 10
	}
	panic(fmt.Errorf("Cannot convert value: %v", n))
}

// Options
func NoteInstr(in int) func(*note) {
	return func(n *note) {
		n.instr = newInt(in)
	}
}

func NoteOctave(o int) func(*note) {
	return func(n *note) {
		n.octave = newInt(o)
	}
}

func NoteNote(o rune) func(*note) {
	return func(n *note) {
		n.note = newRune(o)
	}
}

func NoteAmp(o int) func(*note) {
	return func(n *note) {
		n.amp = newInt(o)
	}
}

func NoteDur(o int) func(*note) {
	return func(n *note) {
		n.dur = newInt(o)
	}
}
