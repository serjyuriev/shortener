package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/serjyuriev/shortener/internal/pkg/handlers"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type Server interface {
	Start() error
}

type server struct {
	address string
	baseURL string
}

func NewServer(address, baseURL, fileStoragePath string) (*server, error) {
	var err error
	handlers.ShortURLHost = baseURL
	handlers.Store, err = storage.NewStore(fileStoragePath)
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
	r.Use(middleware.Recoverer)
	r.Get("/{shortPath}", handlers.GetURLHandler)
	r.Post("/", handlers.PostURLHandler)
	r.Post("/api/shorten", handlers.PostURLApiHandler)
	log.Printf("starting server on %s\n", s.address)
	return http.ListenAndServe(s.address, r)
}
