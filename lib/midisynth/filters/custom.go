package filters

import (
	"fmt"
	"log/slog"
	"math"
	"regexp"

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

var useVRe = regexp.MustCompile(`\bv\b`)

func (cf *Custom) Apply(w waves.Wave) waves.Wave {
	env := cf.prepEnv(0.0, 0.0, nil, nil)
	program, err := expr.Compile(cf.Code, expr.Env(env))
	if err != nil {
		panic(fmt.Errorf("Failed to compile custom code: %w\nCode: `%v`", err, cf.Code))
	}
	return &customImpl{
		wave:    w,
		parent:  cf,
		program: program,
		useV:    useVRe.MatchString(cf.Code),
	}
}

func (cf *Custom) prepEnv(v float64, t float64, in func(t float64, ctx *waves.NoteCtx) float64, ctx *waves.NoteCtx) map[string]any {
	return map[string]any{
		"v":     v,
		"t":     t,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"input": in,
		"ctx":   ctx,
	}

}

type customImpl struct {
	wave    waves.Wave
	parent  *Custom
	program *vm.Program
	useV    bool
}

func (i *customImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	inputValue := 0.0
	if i.useV {
		inputValue = i.wave.Value(t, ctx)
	}
	env := i.parent.prepEnv(inputValue, t, i.wave.Value, ctx)
	value, err := expr.Run(i.program, env)
	if err != nil {
		slog.Error("Custom expression failed", "error", err)
		return 0
	}
	f, ok := value.(float64)
	if !ok {
		slog.Error("Custom expression incorrect result", "value", value)
		return 0
	}
	return f
}
