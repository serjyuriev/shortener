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

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/service"
	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type userURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type (
	postShortenRequest struct {
		URL string `json:"url"`
	}

	postShortenResponse struct {
		Result string `json:"result"`
	}
)

type (
	postBatchSingleRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	postBatchSingleResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

type ContextKey string

var contextKeyUID = ContextKey("uid")

// Handlers store link to service layer and app's base URL.
type Handlers struct {
	svc     service.Service
	baseURL string
}

// MakeHandlers initializes application handler functions and service layer.
func MakeHandlers() (*Handlers, error) {
	svc, err := service.NewService()
	if err != nil {
		return nil, fmt.Errorf("unable to create new service:\n%w", err)
	}

	cfg := config.GetConfig()
	return &Handlers{
		baseURL: cfg.BaseURL,
		svc:     svc,
	}, nil
}

// DeleteURLsHandler removes URLs provided by user from storage.
func (h *Handlers) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(contextKeyUID).(string)

	var req []string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("unable to decode request's body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if len(req) == 0 {
		http.Error(w, "Body cannot be empty.", http.StatusBadRequest)
		return
	}

	h.svc.DeleteURLs(uid, req)
	w.WriteHeader(http.StatusAccepted)
}

// GetURLHandler searches service store for provided short URL
// and, if such URL is found, sends a response,
// redirecting to the corresponding long URL.
func (h *Handlers) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortPath := strings.TrimPrefix(r.URL.Path, "/")
	if shortPath == "" {
		http.Error(w, "No short URL is provided.", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	original, err := h.svc.FindOriginalURL(ctx, shortPath)
	if err != nil {
		if errors.Is(err, storage.ErrShortenedDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		log.Printf("unable to find full URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", string(original))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// GetUserURLsAPIHandler returns all URLs that were added by current user.
func (h *Handlers) GetUserURLsAPIHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(contextKeyUID).(string)
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	m, err := h.svc.FindURLsByUser(ctx, uid)
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
	res := make([]userURLs, 0)
	for key, val := range m {
		res = append(res, userURLs{
			ShortURL:    fmt.Sprintf("%s/%s", h.baseURL, key),
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

// PingHandler provides health status of application.
func (h *Handlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := h.svc.Ping(ctx); err != nil {
		log.Printf("unable to ping database: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

// PostBatchHandler adds URLs provided by user into storage,
// returning shortened URLs with corresponding correlation ID.
func (h *Handlers) PostBatchHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(contextKeyUID).(string)
	var req []postBatchSingleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("unable to decode request's body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	res := make([]postBatchSingleResponse, 0)
	m := make(map[string]string)
	for _, sreq := range req {
		s := shorty.GenerateShortPath()
		sres := postBatchSingleResponse{
			CorrelationID: sreq.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", h.baseURL, s),
		}
		res = append(res, sres)
		m[s] = sreq.OriginalURL
	}

	if err := h.svc.InsertManyURLs(r.Context(), uid, m); err != nil {
		log.Printf("unable to insert urls: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
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

// PostURLApiHandler adds single URL provided by user in JSON format into storage,
// returning its generated short URL.
func (h *Handlers) PostURLApiHandler(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	hadConflict := false
	if err := h.svc.InsertNewURLPair(ctx, uid, s, req.URL); err != nil {
		if !errors.Is(err, storage.ErrNotUniqueOriginalURL) {
			log.Printf("unable to save URL: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		s, err = h.svc.FindByOriginalURL(r.Context(), req.URL)
		if err != nil {
			log.Printf("unable to find original URL: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		hadConflict = true
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, s)
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
	if hadConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(json)
}

// PostURLHandler reads a long URL provided in request body
// and, if successful, creates a corresponding short URL,
// storing both in service store.
func (h *Handlers) PostURLHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err = url.ParseRequestURI(string(b)); err != nil {
		log.Printf("unable to parse request URL: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	s := shorty.GenerateShortPath()
	uid := r.Context().Value(contextKeyUID).(string)
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	hadConflict := false
	if err = h.svc.InsertNewURLPair(ctx, uid, s, string(b)); err != nil {
		if !errors.Is(err, storage.ErrNotUniqueOriginalURL) {
			log.Printf("unable to save URL: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		s, err = h.svc.FindByOriginalURL(r.Context(), string(b))
		if err != nil {
			log.Printf("unable to find original URL: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		hadConflict = true
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, s)
	w.Header().Set("Content-Type", "text/plain")
	if hadConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(shortURL))
}
