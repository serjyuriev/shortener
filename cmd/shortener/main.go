package main

import (
	"log"

	"github.com/serjyuriev/shortener/internal/pkg/server"
)

func main() {
	s := server.NewServer("localhost", 8080)
	log.Fatal(s.Start())
}
