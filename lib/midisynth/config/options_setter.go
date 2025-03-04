package config

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type param struct {
	name  string
	value any
}

func Param(name string, value any) *param {
	return &param{
		name:  name,
		value: value,
	}
}

func SetOptions(obj, options any, params ...*param) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetOptions failed: %v\n%s", r, debug.Stack())
		}
	}()

	// Set filter options
	switch opts := options.(type) {
	case float64, int, string:
		if err := setSingleOption(obj, opts); err != nil {
			return err
		}
	case map[string]any:
		if err := setMapOptions(obj, opts); err != nil {
			return err
		}
	case nil:
		// Do nothing
	default:
		return fmt.Errorf("Unsupported options type: %v (%T)", options, options)
	}

	// Set global params
	for _, p := range params {
		setParamByTagName(obj, p.name, p.value)
	}
	return
}

func setSingleOption(obj any, f any) error {
	v := reflect.ValueOf(obj).Elem()
	if v.NumField() != 1 {
		return fmt.Errorf("Struct should contain exactly 1 field, found: %+v", v.Interface())
	}
	fld := v.Field(0)
	assignToField(fld, f)
	return nil
}

func assignToField(fld reflect.Value, f any) {
	fldType := fld.Type().String()
	switch fldType {
	case "float64":
		switch x := f.(type) {
		case int:
			fld.SetFloat(float64(x))
			return
		case int64:
			fld.SetFloat(float64(x))
			return
		}
	case "int", "int64":
		switch x := f.(type) {
		case int:
			fld.SetInt(int64(x))
			return
		case int64:
			fld.SetInt(x)
			return
		}
	case "waves.Wave":
		if name, ok := f.(string); ok {
			if wave, ok := waves.Waves[name]; ok {
				fld.Set(reflect.ValueOf(wave))
				return
			}
			panic(fmt.Errorf("Unknown wave: %v", name))
		}
	}

	fld.Set(reflect.ValueOf(f))

}

func setMapOptions(obj any, opts map[string]any) error {
	for key, value := range opts {
		err := setOptionByName(obj, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func setOptionByName(obj any, name string, f any) error {
	typ := reflect.TypeOf(obj).Elem()
	n := typ.NumField()
	knownFields := []string{}
	for i := range n {
		fld := typ.Field(i)
		tagsRaw := fld.Tag.Get("option")
		tags := strings.Split(tagsRaw, ",")

		// Check field name
		// Search by tags
		if strings.ToLower(fld.Name) == name || slices.Contains(tags, name) {
			v := reflect.ValueOf(obj).Elem()
			assignToField(v.Field(i), f)
			return nil
		}

		knownFields = append(knownFields, fld.Name)
	}
	return fmt.Errorf("Cannot save option '%v' (%v) into %+v (fields: %v)", name, f, obj, knownFields)
}

func setParamByTagName(obj any, name string, f any) {
	typ := reflect.TypeOf(obj).Elem()
	n := typ.NumField()
	for i := range n {
		fld := typ.Field(i)
		tag := fld.Tag.Get("param")

		// Search by tags
		if name == tag {
			v := reflect.ValueOf(obj).Elem()
			assignToField(v.Field(i), f)
			return
		}
	}
	// Do nothing
}
