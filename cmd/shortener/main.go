package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"

	"github.com/serjyuriev/shortener/internal/pkg/server"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var cfg *config

func init() {
	cfg = &config{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "web server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL for shorten links")
	flag.StringVar(&cfg.FileStoragePath, "f", "shorten.json", "shorten URL file path")
}

func main() {
	flag.Parse()
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(cfg.BaseURL, cfg.ServerAddress) {
		cfg.BaseURL = fmt.Sprintf("http://%s", cfg.ServerAddress)
	}
	log.Println(cfg)
	s, err := server.NewServer(cfg.ServerAddress, cfg.BaseURL, cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start())
}
