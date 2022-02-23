package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/serjyuriev/shortener/internal/pkg/server"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	s, err := server.NewServer(cfg.ServerAddress, cfg.BaseURL, cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start())
}
