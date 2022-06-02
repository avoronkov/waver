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
