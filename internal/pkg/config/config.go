package config

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Protocol        string `env:"-"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
}

var once sync.Once
var cfg *Config

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL for shorten links")
		flag.StringVar(&cfg.DatabaseDSN, "d", "", "data source name")
		flag.StringVar(&cfg.FileStoragePath, "f", "shorten.json", "shorten URL file path")
		flag.StringVar(&cfg.Protocol, "p", "http", "protocol to use (http/https)")
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "web server address")
		flag.Parse()

		if err := env.Parse(cfg); err != nil {
			log.Fatalf("unable to load values from environment variables: %v", err)
		}
		if !strings.Contains(cfg.BaseURL, cfg.ServerAddress) {
			cfg.BaseURL = fmt.Sprintf("%s://%s", cfg.Protocol, cfg.ServerAddress)
		}
		log.Println(cfg)
	})
	return cfg
}
