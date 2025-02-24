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
		dur:   Var("_dur"),
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

var _ types.Signaler = (*note)(nil)

func (n *note) Eval(bit int64, ctx types.Context) (res []signals.Signal) {
	instr := n.instr.Val(bit, ctx)
	if s, ok := instr.(Str); ok {
		return n.evalInstr(bit, ctx, string(s))
	} else if l, ok := instr.(EvaluatedList); ok {
		for i := range l.Len() {
			item := l.Get(i)
			s := item.(Str)
			res = append(res, n.evalInstr(bit, ctx, string(s))...)
		}
		return res
	}
	panic(fmt.Errorf("Don't know how to use for instrument: %v (%T)", instr, instr))
}

func (n *note) evalInstr(bit int64, ctx types.Context, in string) (res []signals.Signal) {
	nt := n.note.Val(bit, ctx)
	if v, ok := nt.(Num); ok {
		nt := n.seqNoteNumberToNote(int64(v))
		return n.evalInstrNote(bit, ctx, in, nt)
	} else if f, ok := nt.(Float); ok {
		nt := notes.Note{Freq: float64(f)}
		return n.evalInstrNote(bit, ctx, in, nt)
	} else if s, ok := nt.(Str); ok {
		if string(s) == "_" {
			return n.evalInstrNote(bit, ctx, in, notes.Note{})
		}
		panic(fmt.Errorf("Don't know how to use for note: %v (%T)", s, s))
	} else if l, ok := nt.(EvaluatedList); ok {
		for i := range l.Len() {
			item := l.Get(i)
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

func (n *note) evalInstrNote(bit int64, ctx types.Context, in string, nt notes.Note) (res []signals.Signal) {
	amp := n.amp.Val(bit, ctx)
	if v, ok := amp.(Num); ok {
		a := float64(v) / 16.0
		return n.evalIntrNoteAmp(bit, ctx, in, nt, a)
	} else if f, ok := amp.(Float); ok {
		return n.evalIntrNoteAmp(bit, ctx, in, nt, float64(f))
	} else if l, ok := amp.(EvaluatedList); ok {
		for i := range l.Len() {
			item := l.Get(i)
			if x, ok := item.(Num); ok {
				a := float64(x) / 16.0
				res = append(res, n.evalIntrNoteAmp(bit, ctx, in, nt, a)...)
			} else if f, ok := item.(Float); ok {
				res = append(res, n.evalIntrNoteAmp(bit, ctx, in, nt, float64(f))...)
			} else {
				panic(fmt.Errorf("Don't know how to convert to Amplitude: %v (%T)", item, item))
			}
		}
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

func (n *note) evalIntrNoteAmp(bit int64, ctx types.Context, inst string, note notes.Note, amp float64) (res []signals.Signal) {
	dur := n.dur.Val(bit, ctx)
	durList := toInt64List(dur, bit, ctx)
	for _, d := range durList {
		res = append(res, n.format(inst, note, amp, d))
	}
	return
}

func (n *note) format(inst string, note notes.Note, amp float64, dur int64) signals.Signal {
	sig := &signals.Signal{
		Instrument:   inst,
		Note:         note,
		DurationBits: int(dur),
		Amp:          amp,
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

func toInt64List(v types.Value, _ int64, _ types.Context) []int64 {
	if n, ok := v.(Num); ok {
		return []int64{int64(n)}
	}
	if l, ok := v.(EvaluatedList); ok {
		var res []int64
		for i := range l.Len() {
			item := l.Get(i)
			i := item.(Num)
			res = append(res, int64(i))
		}
		return res
	}
	panic(fmt.Errorf("Don't know how to convert to []int64: %v", v))
}
