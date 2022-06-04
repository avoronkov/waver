package main

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq"
)

func newInt(n int) *int {
	return &n
}

func newRune(r rune) *rune {
	return &r
}

type note struct {
	instr ValueFn
	note  ValueFn
	amp   ValueFn
	dur   ValueFn
}

func Note(instr, nt ValueFn, opts ...func(*note)) seq.Signaler {
	n := &note{
		instr: instr,
		note:  nt,
		amp:   Const(5),
		dur:   Const(4),
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

var _ seq.Signaler = (*note)(nil)

func (n *note) Eval(bit int64, ctx seq.Context) (res []string) {
	instr := n.instr.Val(bit, ctx)
	for _, i := range instr.ToInt64List() {
		res = append(res, n.evalInstr(bit, ctx, i)...)
	}
	return
}

func (n *note) evalInstr(bit int64, ctx seq.Context, in int64) (res []string) {
	nt := n.note.Val(bit, ctx)
	for _, i := range nt.ToInt64List() {
		res = append(res, n.evalInstrNote(bit, ctx, in, i)...)
	}
	return
}

func (n *note) evalInstrNote(bit int64, ctx seq.Context, in, nt int64) (res []string) {
	amp := n.amp.Val(bit, ctx)
	for _, i := range amp.ToInt64List() {
		res = append(res, n.evalIntrNoteAmp(bit, ctx, in, nt, i)...)
	}
	return
}

func (n *note) evalIntrNoteAmp(bit int64, ctx seq.Context, inst, note, amp int64) (res []string) {
	dur := n.dur.Val(bit, ctx)
	durList := dur.ToInt64List()
	for _, d := range durList {
		res = append(res, n.format(inst, note, amp, d))
	}
	return
}

func (n *note) format(inst, note, amp, dur int64) string {
	return fmt.Sprintf("%c%s%c%c", toRune(inst), KeyNote(note).String(), toRune(amp), toRune(dur))
}

func toRune(n int64) rune {
	if n >= 0 && n < 10 {
		return '0' + rune(n)
	}
	if n >= 10 {
		return 'A' + rune(n) - 10
	}
	panic(fmt.Errorf("Cannot convert value: %v", n))
}

// Options
func NoteInstr(in ValueFn) func(*note) {
	return func(n *note) {
		n.instr = in
	}
}

func NoteNote(o ValueFn) func(*note) {
	return func(n *note) {
		n.note = o
	}
}

func NoteAmp(o ValueFn) func(*note) {
	return func(n *note) {
		n.amp = o
	}
}

func NoteDur(o ValueFn) func(*note) {
	return func(n *note) {
		n.dur = o
	}
}
