package main

import (
	"log"

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/server"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	cfg := config.GetConfig()
	log.Println(cfg)

	var s server.Server
	var err error
	if cfg.DatabaseDSN == "" {
		s, err = server.NewServer(cfg.ServerAddress, cfg.BaseURL, cfg.FileStoragePath, false)
	} else {
		s, err = server.NewServer(cfg.ServerAddress, cfg.BaseURL, cfg.DatabaseDSN, true)
	}
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
	log.Fatal(s.Start())
}
