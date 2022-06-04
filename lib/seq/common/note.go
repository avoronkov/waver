package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/types"
)

type note struct {
	scale notes.Scale
	instr types.ValueFn
	note  types.ValueFn
	amp   types.ValueFn
	dur   types.ValueFn
}

func Note(scale notes.Scale, instr, nt types.ValueFn, opts ...func(*note)) types.Signaler {
	n := &note{
		scale: scale,
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
	for _, i := range toInt64List(instr, bit, ctx) {
		res = append(res, n.evalInstr(bit, ctx, i)...)
	}
	return
}

func (n *note) evalInstr(bit int64, ctx types.Context, in int64) (res []signals.Signal) {
	nt := n.note.Val(bit, ctx)
	if v, ok := nt.(Num); ok {
		nt := n.seqNoteNumberToNote(int64(v))
		return n.evalInstrNote(bit, ctx, in, nt)
	} else if f, ok := nt.(Float); ok {
		nt := notes.Note{Freq: float64(f)}
		return n.evalInstrNote(bit, ctx, in, nt)
	} else if l, ok := nt.(List); ok {
		llen := l.Len()
		for i := 0; i < llen; i++ {
			item := l.Get(i, bit, ctx)
			if i, ok := item.(Num); ok {
				nt := n.seqNoteNumberToNote(int64(i))
				res = append(res, n.evalInstrNote(bit, ctx, in, nt)...)
			} else if f, ok := item.(Float); ok {
				nt := notes.Note{Freq: float64(f)}
				res = append(res, n.evalInstrNote(bit, ctx, in, nt)...)
			} else {
				panic(fmt.Errorf("Don't know how to convert to Note: %v (%T)", item, item))
			}
		}
	}
	return
}

func (n *note) evalInstrNote(bit int64, ctx types.Context, in int64, nt notes.Note) (res []signals.Signal) {
	amp := n.amp.Val(bit, ctx)
	for _, i := range toInt64List(amp, bit, ctx) {
		res = append(res, n.evalIntrNoteAmp(bit, ctx, in, nt, i)...)
	}
	return
}

func (n *note) seqNoteNumberToNote(num int64) notes.Note {
	nt, ok := n.scale.ByNumber(int(num))
	if !ok {
		panic(fmt.Errorf("Unknown note index: %v", num))
	}
	return nt
}

func (n *note) evalIntrNoteAmp(bit int64, ctx types.Context, inst int64, note notes.Note, amp int64) (res []signals.Signal) {
	dur := n.dur.Val(bit, ctx)
	durList := toInt64List(dur, bit, ctx)
	for _, d := range durList {
		res = append(res, n.format(inst, note, amp, d))
	}
	return
}

func (n *note) format(inst int64, note notes.Note, amp, dur int64) signals.Signal {
	sig := &signals.Signal{
		Instrument:   int(inst),
		Note:         note,
		DurationBits: int(dur),
		Amp:          float64(amp) / 16.0,
	}
	return *sig
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

func toInt64List(v types.Value, bit int64, ctx types.Context) []int64 {
	if n, ok := v.(Num); ok {
		return []int64{int64(n)}
	}
	if l, ok := v.(List); ok {
		var res []int64
		llen := l.Len()
		for i := 0; i < llen; i++ {
			item := l.Get(i, bit, ctx)
			i := item.(Num)
			res = append(res, int64(i))
		}
		return res
	}
	panic(fmt.Errorf("Don't know how to convert to []int64: %v", v))
}
