// Package app implements a service
// for shortening URLs.
package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	UrlLength int
	Host      string
	Port      int
}

var keyValueStore map[string]string
var config Config
var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// MakeService initializes shortener service params
// and returns a ServeMux with handlers.
func MakeService(cfg Config) *http.ServeMux {
	rand.Seed(time.Now().UnixNano())
	keyValueStore = make(map[string]string, 0)
	config = cfg

	mux := &http.ServeMux{}
	mux.HandleFunc("/", rootHandler)

	return mux
}

// rootHandler calls specific handlers,
// based on the HTTP request method.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUrlHandler(w, r)
	case http.MethodPost:
		postUrlHandler(w, r)
	default:
		http.Error(w, "Only GET and POST requests are allowed.",
			http.StatusMethodNotAllowed)
		return
	}
}

// getUrlHandler searches service store for provided short URL
// and, if such URL is found, sends a response,
// redirecting to the corresponding long URL.
func getUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, "No short URL is provided.", http.StatusBadRequest)
		return
	}
	shortPath := r.URL.Path[1:]
	longUrl, ok := keyValueStore[shortPath]
	if !ok {
		http.Error(w, "No URL was found.", http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", longUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// postUrlHandler reads a long URL provided in request body
// and, if successful, creates a corresponding short URL,
// storing both in service store.
func postUrlHandler(w http.ResponseWriter, r *http.Request) {
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

	s := processLongUrl(string(b))
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(s))
}

// processLongUrl stores a key-value pair of short URL path and long URL
// in service store and returns short URL.
func processLongUrl(longUrl string) string {
	shortPath := generateShortPath()
	keyValueStore[shortPath] = longUrl
	return fmt.Sprintf("%s:%d/%s", config.Host, config.Port, shortPath)
}

// generateShortPath generates a pseudorandom
// letter sequence of fixed length.
func generateShortPath() string {
	b := make([]rune, config.UrlLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
