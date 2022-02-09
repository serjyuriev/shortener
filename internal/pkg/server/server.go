package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

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
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/{shortPath}", handlers.GetURLHandler)
	r.Post("/", handlers.PostURLHandler)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), r)
}
