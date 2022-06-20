package config

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config contains information about application configuration.
type Config struct {
	BaseURL         string `env:"BASE_URL" envDefault:"https://localhost:8080"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Protocol        string `env:"-"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" envDefault:"true"`
}

// String prints current configuration.
func (c *Config) String() string {
	return fmt.Sprintf(`

	loaded configuration
		BaseURL:         %s
		DatabaseDSN:     %s
		FileStoragePath: %s
		Protocol:        %s
		ServerAddress:   %s
	`, c.BaseURL, c.DatabaseDSN, c.FileStoragePath, c.Protocol, c.ServerAddress)
}

var once sync.Once
var cfg *Config

// GetConfig parses flags and environment variables and creates current configuration object.
// Environment variables will override values provided by flags.
func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		flag.StringVar(&cfg.BaseURL, "b", "https://localhost:8080", "base URL for shorten links")
		flag.StringVar(&cfg.DatabaseDSN, "d", "", "data source name")
		flag.StringVar(&cfg.FileStoragePath, "f", "shorten.json", "shorten URL file path")
		flag.StringVar(&cfg.Protocol, "p", "https", "protocol to use (http/https)")
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "web server address")
		flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable https")
		flag.Parse()

		if err := env.Parse(cfg); err != nil {
			log.Fatalf("unable to load values from environment variables: %v", err)
		}
		if !strings.Contains(cfg.BaseURL, cfg.ServerAddress) {
			cfg.BaseURL = fmt.Sprintf("%s://%s", cfg.Protocol, cfg.ServerAddress)
		}
	})
	return cfg
}
