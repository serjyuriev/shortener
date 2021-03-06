package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/serjyuriev/shortener/internal/pkg/handlers"
)

var errInvalidCookie = errors.New("can not validate signature")
var cookieName = "userID"
var key = []byte("sh0rt7")

var contextKeyUID = handlers.ContextKey("uid")

// Auth checks whether current request has cookie and either
// generates new cookie if no cookie was provided or
// validates provided cookie.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("error getting %s cookie from the request: %v", cookieName, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			uid, newCookieValue := generateNewUserIDCookie()
			newCookie := &http.Cookie{
				Name:    cookieName,
				Value:   newCookieValue,
				Expires: time.Now().Add(600 * time.Second).UTC(),
				Path:    "/",
			}
			http.SetCookie(w, newCookie)
			ctx := context.WithValue(r.Context(), contextKeyUID, uid.String())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		uid, err := validateCookie(cookie.Value)
		if err != nil {
			if !errors.Is(err, errInvalidCookie) {
				log.Printf("unable to validate cookie: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			var newCookieValue string
			uid, newCookieValue = generateNewUserIDCookie()
			newCookie := &http.Cookie{
				Name:    cookieName,
				Value:   newCookieValue,
				Expires: time.Now().Add(60 * time.Second).UTC(),
			}
			http.SetCookie(w, newCookie)
		}
		ctx := context.WithValue(r.Context(), contextKeyUID, uid.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Gzipper checks whether current request was encoded with gzip
// and if so decodes it.
func Gzipper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("unable to create gzip reader: %v\n", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			defer gr.Close()

			body, err := io.ReadAll(gr)
			if err != nil {
				log.Printf("unable to read request body: %v\n", err)
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r.ContentLength = int64(len(body))
		}
		next.ServeHTTP(w, r)
	})
}

func generateNewUserIDCookie() (uuid.UUID, string) {
	uid := uuid.New()
	h := hmac.New(sha256.New, key)
	h.Write([]byte(uid.String()))
	cookie := append([]byte(uid.String()), h.Sum(nil)...)
	return uid, hex.EncodeToString(cookie)
}

func validateCookie(cookie string) (uuid.UUID, error) {
	decodedCookie, err := hex.DecodeString(cookie)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to decode cookie:\n%w", err)
	}

	h := hmac.New(sha256.New, key)
	if _, err = h.Write(decodedCookie[:36]); err != nil {
		return uuid.Nil, fmt.Errorf("unable to write decoded cookie:\n%w", err)
	}

	if hmac.Equal(decodedCookie[36:], h.Sum(nil)) {
		uid, err := uuid.Parse(string(decodedCookie[:36]))
		if err != nil {
			return uuid.Nil, fmt.Errorf("unable to parse uuid:\n%w", err)
		}
		return uid, nil
	}
	return uuid.Nil, errInvalidCookie
}
