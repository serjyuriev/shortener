package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/serjyuriev/shortener/internal/pkg/handlers"
)

type Server interface {
	Start() error
}

type server struct {
	host string
	port int
}

func NewServer(host string, port int) *server {
	return &server{
		host: host,
		port: port,
	}
}

func (s *server) Start() error {
	h, err := handlers.MakeHandlers(context.Background(), s.host, s.port)
	if err != nil {
		return err
	}
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), h)
}
