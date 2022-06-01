package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Go Web Assembly")

	js.Global().Set("goPlay", js.FuncOf(goPlay))
	<-c
}
