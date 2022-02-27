package common

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/midisynth/udp"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type note struct {
	instr types.ValueFn
	note  types.ValueFn
	amp   types.ValueFn
	dur   types.ValueFn
}

func Note(instr, nt types.ValueFn, opts ...func(*note)) types.Signaler {
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

var _ types.Signaler = (*note)(nil)

func (n *note) Eval(bit int64, ctx types.Context) (res []signals.Signal) {
	instr := n.instr.Val(bit, ctx)
	for _, i := range toInt64List(instr) {
		res = append(res, n.evalInstr(bit, ctx, i)...)
	}
	return
}

func (n *note) evalInstr(bit int64, ctx types.Context, in int64) (res []signals.Signal) {
	nt := n.note.Val(bit, ctx)
	for _, i := range toInt64List(nt) {
		res = append(res, n.evalInstrNote(bit, ctx, in, i)...)
	}
	return
}

func (n *note) evalInstrNote(bit int64, ctx types.Context, in, nt int64) (res []signals.Signal) {
	amp := n.amp.Val(bit, ctx)
	for _, i := range toInt64List(amp) {
		res = append(res, n.evalIntrNoteAmp(bit, ctx, in, nt, i)...)
	}
	return
}

func (n *note) evalIntrNoteAmp(bit int64, ctx types.Context, inst, note, amp int64) (res []signals.Signal) {
	dur := n.dur.Val(bit, ctx)
	durList := toInt64List(dur)
	for _, d := range durList {
		res = append(res, n.format(inst, note, amp, d))
	}
	return
}

func (n *note) format(inst, note, amp, dur int64) signals.Signal {
	noteT, err := ParseStandardNoteCode(note)
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%c%s%c%c", toRune(inst), noteT.String(), toRune(amp), toRune(dur))
	sig, err := udp.ParseMessage([]byte(s))
	if err != nil {
		panic(err)
	}
	return *sig
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
func NoteInstr(in types.ValueFn) func(*note) {
	return func(n *note) {
		n.instr = in
	}
}

func NoteNote(o types.ValueFn) func(*note) {
	return func(n *note) {
		n.note = o
	}
}

func NoteAmp(o types.ValueFn) func(*note) {
	return func(n *note) {
		n.amp = o
	}
}

func NoteDur(o types.ValueFn) func(*note) {
	return func(n *note) {
		n.dur = o
	}
}

func toInt64List(v types.Value) []int64 {
	if n, ok := v.(Num); ok {
		return []int64{int64(n)}
	}
	if l, ok := v.(List); ok {
		var res []int64
		for _, item := range l {
			i := item.(Num)
			res = append(res, int64(i))
		}
		return res
	}
	panic(fmt.Errorf("Don't know how to convert to []int64: %v", v))
}
