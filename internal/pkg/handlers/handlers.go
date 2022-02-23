package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type postShortenRequest struct {
	URL string `json:"url"`
}

type postShortenResponse struct {
	Result string `json:"result"`
}

var Store storage.Store
var ShortURLHost string

func PostURLApiHandler(w http.ResponseWriter, r *http.Request) {
	var req postShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.URL) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := shorty.GenerateShortPath()
	if err := Store.InsertNewURLPair(storage.ShortPath(s), storage.LongURL(req.URL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shortURL := fmt.Sprintf("%s/%s", ShortURLHost, s)
	res := postShortenResponse{
		Result: shortURL,
	}
	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

// PostURLHandler reads a long URL provided in request body
// and, if successful, creates a corresponding short URL,
// storing both in service store.
func PostURLHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(string(b)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := shorty.GenerateShortPath()
	err = Store.InsertNewURLPair(storage.ShortPath(s), storage.LongURL(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shortURL := fmt.Sprintf("%s/%s", ShortURLHost, s)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// GetURLHandler searches service store for provided short URL
// and, if such URL is found, sends a response,
// redirecting to the corresponding long URL.
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortPath := strings.TrimPrefix(r.URL.Path, "/")
	if shortPath == "" {
		http.Error(w, "No short URL is provided.", http.StatusBadRequest)
		return
	}
	l, err := Store.FindLongURL(storage.ShortPath(shortPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", string(l))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
