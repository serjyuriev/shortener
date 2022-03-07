package middleware

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"net/http"
)

func Gzipper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("unable to create gzip reader: %v\n", err)
				io.WriteString(w, err.Error())
				return
			}
			defer gr.Close()

			body, err := io.ReadAll(gr)
			if err != nil {
				if !errors.Is(err, io.ErrUnexpectedEOF) &&
					!errors.Is(err, io.EOF) {
					log.Printf("unable to read request body: %v\n", err)
					io.WriteString(w, err.Error())
					return
				}
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r.ContentLength = int64(len(body))
		}
		next.ServeHTTP(w, r)
	})
}
