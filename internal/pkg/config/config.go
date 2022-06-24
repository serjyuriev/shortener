package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config contains information about application configuration.
type Config struct {
	ConfigPath      string `json:"-" env:"CONFIG"`
	BaseURL         string `json:"base_url" env:"BASE_URL" envDefault:"https://localhost:8080"`
	DatabaseDSN     string `json:"database_dsn,omitempty" env:"DATABASE_DSN"`
	FileStoragePath string `json:"file_storage_path,omitempty" env:"FILE_STORAGE_PATH"`
	Protocol        string `json:"protocol" env:"-"`
	ServerAddress   string `json:"server_address" env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	EnableHTTPS     bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
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
		var configPath string

		flag.StringVar(&configPath, "c", "", "json config file")
		flag.StringVar(&cfg.BaseURL, "b", "https://localhost:8080", "base URL for shorten links")
		flag.StringVar(&cfg.DatabaseDSN, "d", "", "data source name")
		flag.StringVar(&cfg.FileStoragePath, "f", "shorten.json", "shorten URL file path")
		flag.StringVar(&cfg.Protocol, "p", "https", "protocol to use (http/https)")
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "web server address")
		flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable https")
		flag.Parse()

		if configPath != "" {
			file, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
			if err != nil {
				log.Printf("unable to open JSON config: %v\n", err)
			} else {
				jsonBody, err := io.ReadAll(file)
				if err != nil {
					log.Printf("unable to read JSON config: %v\n", err)
				} else if err = json.Unmarshal(jsonBody, cfg); err != nil {
					log.Printf("unable to unmarshal JSON config: %v\n", err)
				} else {
					log.Println("JSON config file was read successfully")
				}
			}
		}

		if err := env.Parse(cfg); err != nil {
			log.Fatalf("unable to load values from environment variables: %v", err)
		}
		if !strings.Contains(cfg.BaseURL, cfg.ServerAddress) {
			cfg.BaseURL = fmt.Sprintf("%s://%s", cfg.Protocol, cfg.ServerAddress)
		}
	})
	return cfg
}
