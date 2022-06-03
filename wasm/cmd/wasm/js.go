package main

import (
	"log"
	"syscall/js"
)

func jsPlay(this js.Value, inputs []js.Value) any {
	code := inputs[0].String()
	if err := updateCode(code); err != nil {
		log.Printf("Updating code failed: %v", code)
		return 1
	}
	return 0
}

func updateCode(input string) error {
	return goParser.ParseData([]byte(input))
}

func jsGetDefaultCode(this js.Value, inputs []js.Value) any {
	s := `: 4 -> "33A11"
: 5 + 1 -> "34C11"
: 6 + 3 -> "34D11"
: 7 + 4 -> "34E11"
: 8 + 6 -> "34G11"
: 9 + 7 -> "34A11"
: 10 + 8 -> "35C11"
: 11 + 9 -> "35D11"
: 12 + 10 -> "35E11"
: 13 + 11 -> "35G11"`
	return js.ValueOf(s)
}
