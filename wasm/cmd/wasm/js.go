package main

import (
	"log"
	"syscall/js"

	"github.com/avoronkov/waver/etc"
	"github.com/avoronkov/waver/lib/share"
)

func jsPlay(this js.Value, inputs []js.Value) any {
	defer doRecover()

	code := inputs[0].String()
	if err := updateCode(code); err != nil {
		log.Printf("Updating code failed: %v", err)
		return js.ValueOf(err.Error())
	}
	return js.ValueOf("OK")
}

func updateCode(input string) error {
	return goParser.ParseData([]byte(input))
}

func jsGetDefaultCode(this js.Value, inputs []js.Value) any {
	return js.ValueOf(string(etc.DefaultCodeExample))
}

func jsPause(this js.Value, inputs []js.Value) any {
	value := inputs[0].Bool()
	goSequencer.Pause(value)
	return js.ValueOf(0)
}

func jsEncode(this js.Value, inputs []js.Value) any {
	arg := inputs[0].String()
	encoded, err := share.Encode(arg)
	res := map[string]any{}
	if err != nil {
		res["error"] = err.Error()
	} else {
		res["data"] = encoded
	}
	return js.ValueOf(res)
}

func jsDecode(this js.Value, inputs []js.Value) any {
	arg := inputs[0].String()
	decoded, err := share.Decode(arg)
	res := map[string]any{}
	if err != nil {
		res["error"] = err.Error()
	} else {
		res["data"] = decoded
	}
	return js.ValueOf(res)
}
