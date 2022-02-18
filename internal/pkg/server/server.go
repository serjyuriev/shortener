package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/serjyuriev/shortener/internal/pkg/handlers"
)

type Server interface {
	Start() error
}

type server struct {
	address string
	baseURL string
}

func NewServer(address, baseURL string) *server {
	handlers.ShortURLHost = baseURL
	return &server{
		address: address,
		baseURL: baseURL,
	}
}

// Start creates new router, binds handlers and starts http server.
func (s *server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/{shortPath}", handlers.GetURLHandler)
	r.Post("/", handlers.PostURLHandler)
	r.Post("/api/shorten", handlers.PostURLApiHandler)
	return http.ListenAndServe(s.address, r)
}
