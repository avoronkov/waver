package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting web-server at http://localhost:9090/")
	err := http.ListenAndServe(":9090", http.FileServer(http.Dir("./static/assets")))
	if err != nil {
		log.Fatal("Failed to start server", err)
		return
	}
}
