package main

import (
	"log"
	"syscall/js"

	"github.com/avoronkov/waver/etc"
)

func jsPlay(this js.Value, inputs []js.Value) any {
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

func jsUpdateInstruments(this js.Value, inputs []js.Value) any {
	data := inputs[0].String()
	goCfg.UpdateData([]byte(data))
	return js.ValueOf(0)
}

func jsGetDefaultInstruments(this js.Value, inputs []js.Value) any {
	return js.ValueOf(string(etc.DefaultConfig))
}
