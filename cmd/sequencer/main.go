package main

import (
	"fmt"
	"log"
	"net"
)

const port = 49161

func main() {
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
	log.Printf("Listening to UDP on localhost:%v", port)
	for {
		buff := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buff)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			continue
		}
		fmt.Printf("Read: '%s' from %v\n", buff[:n], addr)
	}
}
