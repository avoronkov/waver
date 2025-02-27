package filters

import (
	"testing"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

func BenchmarkCustom(b *testing.B) {
	tw := &testWave{}
	ef := Custom{
		Code: `let e = 2.0; v >= 0 ? v ** e : -((-v)**e)`,
	}
	w := ef.Apply(tw)
	ctx := &waves.NoteCtx{}

	for range b.N {
		w.Value(0.5, ctx)
	}
}
