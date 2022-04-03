package server

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/handlers"
	"github.com/serjyuriev/shortener/internal/pkg/middleware"
)

type Server interface {
	Start() error
}

type server struct {
	address  string
	handlers handlers.Handlers
}

func NewServer() (Server, error) {
	h, err := handlers.MakeHandlers()
	if err != nil {
		return nil, fmt.Errorf("unable to make handlers:\n%w", err)
	}

	cfg := config.GetConfig()
	return &server{
		address:  cfg.ServerAddress,
		handlers: h,
	}, nil
}

// Start creates new router, binds handlers and starts http server.
func (s *server) Start() error {
	r := chi.NewRouter()
	r.Use(chimid.Recoverer)
	r.Use(chimid.Compress(gzip.BestSpeed, zippableTypes...))
	r.Use(middleware.Gzipper)
	r.Use(middleware.Auth)
	r.Delete("/api/user/urls", s.handlers.DeleteURLsHandler)
	r.Get("/ping", s.handlers.PingHandler)
	r.Get("/{shortPath}", s.handlers.GetURLHandler)
	r.Get("/api/user/urls", s.handlers.GetUserURLsAPIHandler)
	r.Post("/", s.handlers.PostURLHandler)
	r.Post("/api/shorten", s.handlers.PostURLApiHandler)
	r.Post("/api/shorten/batch", s.handlers.PostBatchHandler)
	log.Printf("starting server on %s\n", s.address)
	return http.ListenAndServe(s.address, r)
}

var zippableTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}
