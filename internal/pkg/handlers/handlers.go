package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

var store storage.Store

// PostURLHandler reads a long URL provided in request body
// and, if successful, creates a corresponding short URL,
// storing both in service store.
func PostURLHandler(w http.ResponseWriter, r *http.Request) {
	err := makeStoreIfNotExists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	splitted := strings.Split(r.Host, ":")
	shortURL := fmt.Sprintf("http://%s:%s/%s", splitted[0], splitted[1], s)
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// GetURLHandler searches service store for provided short URL
// and, if such URL is found, sends a response,
// redirecting to the corresponding long URL.
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	err := makeStoreIfNotExists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shortPath := r.URL.Path[1:]
	if shortPath == "" {
		http.Error(w, "No short URL is provided.", http.StatusBadRequest)
		return
	}
	l, err := store.FindLongURL(context.Background(), storage.ShortPath(shortPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", string(l))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func makeStoreIfNotExists() error {
	if store == nil {
		var err error
		store, err = storage.NewStore(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
