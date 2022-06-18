package config

import (
	"reflect"
	"testing"
)

type structA struct {
	X int
}

type structB struct {
	Y float64
}

type structC struct {
	X int     `option:"xopt"`
	Y float64 `option:"yopt"`
	Z float64 `option:"zzz,zopt"`
}

func TestSetOptions(t *testing.T) {
	tests := []struct {
		name string
		obj  any
		opts any
		exp  any
	}{
		{"Single int value", &structA{}, 63, &structA{63}},
		{"Single float value", &structB{}, 10.5, &structB{10.5}},
		{"Single int->float value", &structB{}, 10, &structB{10.0}},
		{
			"Complex value by field names",
			&structC{},
			map[string]any{
				"x": 11,
				"y": 12.5,
				"z": 13,
			},
			&structC{11, 12.5, 13.0},
		},
		{
			"Complex value by tags",
			&structC{},
			map[string]any{
				"xopt": 11,
				"yopt": 12.5,
				"zzz":  13,
			},
			&structC{11, 12.5, 13.0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetOptions(test.obj, test.opts)
			if err != nil {
				t.Fatalf("SetOptions failed: %v", err)
			}
			if !reflect.DeepEqual(test.obj, test.exp) {
				t.Fatalf("Incorrect result of SetOpions: want %v, got %v", test.exp, test.obj)
			}
		})
	}
}
