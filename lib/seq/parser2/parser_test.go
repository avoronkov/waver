package parser2

import (
	"strings"
	"testing"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
	"github.com/matryer/is"
)

type TestSeq struct {
	signalers []types.Signaler
	vars map[string]types.ValueFn
}

func (t *TestSeq) Add(s types.Signaler) {
	t.signalers = append(t.signalers, s)
}

func (t *TestSeq) Commit() error {
	return nil
}

func (t *TestSeq) Assign(name string, value types.ValueFn) {
	if t.vars == nil {
		t.vars = map[string]types.ValueFn{}
	}
	t.vars[name] = value
}

func TestSignalStatement(t *testing.T) {
	input := ": 4 -> { kick }"

	seq := &TestSeq{}
	p := New(WithSeq(seq))

	err := p.parseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseReader failed: %v", err)
	}

	if len(seq.signalers) != 1 {
		t.Fatalf("seq.signalers incorrect: %v", seq.signalers)
	}

	signaler := seq.signalers[0]

	exp := [][]signals.Signal{
		{
			{Instrument: "kick", DurationBits: 4, Amp: 0.3125},
		},
		nil,
		nil,
		nil,
	}

	testSignaler(t, signaler, exp)
}

func TestVarAssignment(t *testing.T) {
	input := "x = 12\n"

	seq := &TestSeq{}
	p := New(WithSeq(seq))

	err := p.parseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseReader failed: %v", err)
	}

	if len(seq.vars) != 1 {
		t.Fatalf("seq.vars incorrect: %v", seq.vars)
	}

	v, ok := seq.vars["x"]
	if !ok {
		t.Fatalf("Variable 'x' not found in %v", seq.vars)
	}
	
	val := v.Val(0, types.NewContext())
	if val != common.Num(12) {
		t.Errorf("Incorrect value of variable x: expected Num(12), found %v", val)
	}
}

func TestUserDefinedFunction(t *testing.T) {
	is := is.New(t)

	input := `foo x = 23`

	p := New()

	err := p.parseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseReader failed: %v", err)
	}

	udf, ok := p.userFunctions["foo"]
	if !ok {
		t.Fatalf("User function 'foo' not found")
	}

	is.Equal(udf.Name, "foo")
	is.Equal(udf.Arg, "x")

	val := udf.Fn.Val(0, types.NewContext())
	is.Equal(val, common.Num(23))
}

func testSignaler(t *testing.T, signaler types.Signaler, signals [][]signals.Signal) {
	ctx := types.NewContext()
	for i, exp := range signals {
		act := signaler.Eval(int64(i), ctx)
		if !compareSlices(act, exp, compareSignals) {
			t.Errorf("Incorrect signals on bit %v:\nexpected %v\nactual   %v", i, exp, act)
		}
	}
}

func compareSlices[T any](act, exp []T, cmp func(T, T) bool) bool {
	if len(act) != len(exp) {
		return false
	}
	for i, a := range act {
		e := exp[i]
		if !cmp(a, e) {
			return false
		}
	}
	return true
}

func compareSignals(a, b signals.Signal) bool {
	return a.Note == b.Note && a.Amp == b.Amp && a.DurationBits == b.DurationBits && a.Instrument == b.Instrument && a.Manual == b.Manual
}
