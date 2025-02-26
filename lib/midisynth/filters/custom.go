package filters

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Custom struct {
	Code string `option:"code"`
}

func (Custom) New() Filter {
	return &Custom{}
}

func (cf *Custom) Apply(w waves.Wave) waves.Wave {
	env := map[string]any{
		"v":   float64(0.0),
		"t":   float64(0.0),
		"sin": math.Sin,
		"cos": math.Cos,
	}
	program, err := expr.Compile(cf.Code, expr.Env(env))
	if err != nil {
		panic(fmt.Errorf("Failed to compile custom code: %w\nCode: `%v`", err, cf.Code))
	}
	return &customImpl{
		wave:    w,
		program: program,
	}
}

type customImpl struct {
	wave    waves.Wave
	program *vm.Program
}

func (i *customImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	env := map[string]any{
		"v":   i.wave.Value(t, ctx),
		"t":   t,
		"sin": math.Sin,
		"cos": math.Cos,
	}
	value, err := expr.Run(i.program, env)
	if err != nil {
		slog.Error("Custom expression failed", "error", err)
		return 0
	}
	f, ok := value.(float64)
	if !ok {
		slog.Error("Custom expression incorrect result", "value", value)
		return 0.0
	}
	return f
}
