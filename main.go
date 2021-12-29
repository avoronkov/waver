package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("./main.go")
	if err != nil {
		fatal("Opening input failed", err)
	}
	defer f.Close()

	wav, err := ReadWav(f)
	if err != nil {
		fatal("Error reading WAV file", err)
	}

	fmt.Println(wav)
}

type WithErrorStack interface {
	ErrorStack() string
}

func fatal(message string, err error) {
	var errStack WithErrorStack
	if errors.As(err, &errStack) {
		log.Fatalf("%v: %v", message, errStack.ErrorStack())
	} else {
		log.Fatalf("%v: %v", message, err)
	}
}
