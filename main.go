package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gitlab.com/avoronkov/waver/wav"
)

func main() {
	f, err := os.Open("./cats.wav")
	if err != nil {
		fatal("Opening input failed", err)
	}
	defer f.Close()

	wav, err := wav.ReadWav(f)
	if err != nil {
		fatal("Error reading WAV file", err)
	}

	fmt.Println(wav)

	/*
		out, err := os.OpenFile("output.wav", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fatal("Opening output failed", err)
		}
		defer out.Close()

		if err := wav.Write(out); err != nil {
			fatal("Writing wav failed", err)
		}
	*/

	fmt.Println("OK!")
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
