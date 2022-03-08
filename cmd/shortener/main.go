package main

import (
	"log"

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/server"
)

func main() {
	cfg := config.GetConfig()
	s, err := server.NewServer(cfg.ServerAddress, cfg.BaseURL, cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
	log.Fatal(s.Start())
}
