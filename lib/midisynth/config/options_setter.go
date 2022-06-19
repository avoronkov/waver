package config

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

func SetOptions(obj, options any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetOptions failed: %v\n%s", r, debug.Stack())
		}
	}()

	switch opts := options.(type) {
	case float64, int:
		return setSingleOption(obj, opts)
	case map[string]any:
		return setMapOptions(obj, opts)
	case nil:
		// Do nothing
		return
	default:
		return fmt.Errorf("Unsupported options type: %v (%T)", options, options)
	}
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
	for i := 0; i < n; i++ {
		fld := typ.Field(i)
		tagsRaw := fld.Tag.Get("option")
		tags := strings.Split(tagsRaw, ",")

		// Check field name
		// Search by tags
		if strings.ToLower(fld.Name) == name || contains(name, tags) {
			v := reflect.ValueOf(obj).Elem()
			assignToField(v.Field(i), f)
			return nil
		}
	}
	return fmt.Errorf("Cannot save option '%v' (%v) into %+v", name, f, obj)
}

func contains[T comparable](item T, list []T) bool {
	for _, e := range list {
		if item == e {
			return true
		}
	}
	return false
}
