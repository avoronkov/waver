package main

import (
	"io/fs"
	"log"
	"net/http"

	"wasm/static"
)

func main() {
	fsys, err := fs.Sub(static.Assets, "assets")
	if err != nil {
		log.Fatal("Cannot find directory assets", err)
	}

	log.Printf("Starting web-server at http://localhost:9090/")
	err = http.ListenAndServe(":9090", http.FileServer(http.FS(fsys)))
	if err != nil {
		log.Fatal("Failed to start server", err)
		return
	}
}
