package main

import (
	"log"

	"github.com/serjyuriev/shortener/internal/pkg/server"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
	log.Fatal(server.Start())
}
