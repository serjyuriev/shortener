package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type ContextKey string

var contextKeyUID = ContextKey("uid")

type userURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type getUserURLsResponse []userURLs

type postShortenRequest struct {
	URL string `json:"url"`
}

type postShortenResponse struct {
	Result string `json:"result"`
}

var Store storage.Store
var ShortURLHost string

func GetUserURLsAPIHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(contextKeyUID).(string)
	m, err := Store.FindURLsByUser(uid)
	if err != nil {
		if errors.Is(err, storage.ErrNoURLWasFound) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("unable to find full URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	res := make(getUserURLsResponse, 0)
	for key, val := range m {
		res = append(res, userURLs{
			ShortURL:    fmt.Sprintf("%s/%s", ShortURLHost, key),
			OriginalURL: val,
		})
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Printf("unable to marshal response: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func PostURLApiHandler(w http.ResponseWriter, r *http.Request) {
	var req postShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("unable to decode request's body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(req.URL) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		log.Printf("unable to parse request URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s := shorty.GenerateShortPath()
	uid := r.Context().Value(contextKeyUID).(string)
	if err := Store.InsertNewURLPair(uid, s, req.URL); err != nil {
		log.Printf("unable to save URL: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	shortURL := fmt.Sprintf("%s/%s", ShortURLHost, s)
	res := postShortenResponse{
		Result: shortURL,
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Printf("unable to marshal response: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		log.Printf("unable to read request's body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(b) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(string(b)); err != nil {
		log.Printf("unable to parse request URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	s := shorty.GenerateShortPath()
	uid := r.Context().Value(contextKeyUID).(string)
	err = Store.InsertNewURLPair(uid, s, string(b))
	if err != nil {
		log.Printf("unable to save URL: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
	l, err := Store.FindLongURL(shortPath)
	if err != nil {
		log.Printf("unable to find full URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", string(l))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Store.Ping(ctx); err != nil {
		log.Printf("unable to ping database: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}
