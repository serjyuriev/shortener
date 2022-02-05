package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type serverConfig struct {
	host string
	port int
}

var store storage.Store
var cfg serverConfig

func MakeHandlers(ctx context.Context, host string, port int) (*http.ServeMux, error) {
	cfg = serverConfig{
		host: host,
		port: port,
	}
	var err error
	store, err = storage.NewStore(ctx)
	if err != nil {
		return nil, err
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", rootHandler)

	return mux, nil
}

// rootHandler calls specific handlers,
// based on the HTTP request method.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getURLHandler(w, r)
	case http.MethodPost:
		postURLHandler(w, r)
	default:
		http.Error(w, "Only GET and POST requests are allowed.",
			http.StatusMethodNotAllowed)
		return
	}
}

// postURLHandler reads a long URL provided in request body
// and, if successful, creates a corresponding short URL,
// storing both in service store.
func postURLHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if len(b) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	} else if _, err := url.ParseRequestURI(string(b)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := shorty.GenerateShortPath()
	err = store.InsertNewURLPair(context.Background(), storage.ShortPath(s), storage.LongURL(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shortURL := fmt.Sprintf("http://%s:%d/%s", cfg.host, cfg.port, s)
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// getURLHandler searches service store for provided short URL
// and, if such URL is found, sends a response,
// redirecting to the corresponding long URL.
func getURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "No short URL is provided.", http.StatusBadRequest)
		return
	}
	shortPath := r.URL.Path[1:]
	l, err := store.FindLongURL(context.Background(), storage.ShortPath(shortPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", string(l))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
