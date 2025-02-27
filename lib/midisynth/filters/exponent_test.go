package filters

import (
	"testing"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type testWave struct{}

func (t *testWave) Value(tm float64, ctx *waves.NoteCtx) float64 {
	return 0.25
}

func BenchmarkExponent(b *testing.B) {
	tw := &testWave{}
	ef := Exponent{Value: 2.0}
	w := ef.Apply(tw)
	ctx := &waves.NoteCtx{}

	for range b.N {
		w.Value(0.5, ctx)
	}
}
