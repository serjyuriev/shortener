package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/serjyuriev/shortener/internal/pkg/server"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = "localhost:8080"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	s := server.NewServer(cfg.ServerAddress, cfg.BaseURL)
	log.Fatal(s.Start())
}
