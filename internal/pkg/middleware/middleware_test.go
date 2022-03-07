package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Gzipper(t *testing.T) {
	tests := []struct {
		name               string
		requestURL         string
		testURL            string
		useContentEncoding bool
	}{
		{
			name:               "positive test - with encoding",
			requestURL:         "http://localhost:8080/",
			testURL:            "https://github.com/serjyuriev/",
			useContentEncoding: true,
		},
		{
			name:               "positive test - without encoding",
			requestURL:         "http://localhost:8080/",
			testURL:            "https://github.com/serjyuriev/",
			useContentEncoding: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody bytes.Buffer
			if tt.useContentEncoding {
				w := gzip.NewWriter(&requestBody)

				_, err := w.Write([]byte(tt.testURL))
				require.NoError(t, err)
				err = w.Close()
				require.NoError(t, err)
			} else {
				w := bufio.NewWriter(&requestBody)
				_, err := w.Write([]byte(tt.testURL))
				require.NoError(t, err)
				require.NoError(t, w.Flush())
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				URL, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, tt.testURL, string(URL))
			})
			mid := Gzipper(nextHandler)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, &requestBody)
			if tt.useContentEncoding {
				request.Header.Set("Content-Encoding", "gzip")
			}
			mid.ServeHTTP(recorder, request)
		})
	}
}
