package server

import (
	"compress/gzip"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"

	"github.com/serjyuriev/shortener/internal/pkg/handlers"
	"github.com/serjyuriev/shortener/internal/pkg/middleware"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type Server interface {
	Start() error
}

type server struct {
	address string
	baseURL string
}

func NewServer(address, baseURL, connectionString string, useDB bool) (*server, error) {
	var err error
	handlers.ShortURLHost = baseURL
	if useDB {
		handlers.Store, err = storage.NewPgStore(connectionString)
	} else {
		handlers.Store, err = storage.NewFileStore(connectionString)
	}
	if err != nil {
		return nil, err
	}
	return &server{
		address: address,
		baseURL: baseURL,
	}, nil
}

// Start creates new router, binds handlers and starts http server.
func (s *server) Start() error {
	r := chi.NewRouter()
	r.Use(chimid.Recoverer)
	r.Use(chimid.Compress(gzip.BestSpeed, zippableTypes...))
	r.Use(middleware.Gzipper)
	r.Use(middleware.Auth)
	r.Get("/ping", handlers.PingHandler)
	r.Get("/{shortPath}", handlers.GetURLHandler)
	r.Get("/api/user/urls", handlers.GetUserURLsAPIHandler)
	r.Post("/", handlers.PostURLHandler)
	r.Post("/api/shorten", handlers.PostURLApiHandler)
	r.Post("/api/shorten/batch", handlers.PostBatchHandler)
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
