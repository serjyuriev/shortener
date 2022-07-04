package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Auth(t *testing.T) {
	tests := []struct {
		name       string
		requestURL string
		testURL    string
	}{
		{
			name:       "positive test",
			requestURL: "http://localhost:8080/",
			testURL:    "https://github.com/serjyuriev/",
		},
		{
			name:       "positive test - without encoding",
			requestURL: "http://localhost:8080/",
			testURL:    "https://github.com/serjyuriev/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody bytes.Buffer
			w := bufio.NewWriter(&requestBody)
			_, err := w.Write([]byte(tt.testURL))
			require.NoError(t, err)
			require.NoError(t, w.Flush())

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				URL, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, tt.testURL, string(URL))
			})
			mid := Auth(nextHandler)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, &requestBody)
			request.AddCookie(&http.Cookie{
				Name:  "userID",
				Value: "123",
			})
			mid.ServeHTTP(recorder, request)
		})
	}
}

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

func Test_generateNewUserIDCookie(t *testing.T) {
	uuidOrig, hex := generateNewUserIDCookie()
	uuidParsed, err := validateCookie(hex)
	assert.NoError(t, err)
	assert.Equal(t, uuidOrig, uuidParsed)
}

func Test_validateCookie(t *testing.T) {
	type want struct {
		uid string
	}
	tests := []struct {
		name   string
		cookie string
		want   want
	}{
		{
			name:   "correct cookie",
			cookie: "36353737663139312d613031322d346631362d616665342d36656430643534326535323333d92c2814ac0a09d73a675e9a187324028274fd0b7b03f488db1631c4b5328a",
			want: want{
				uid: "6577f191-a012-4f16-afe4-6ed0d542e523",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, err := validateCookie(tt.cookie)
			require.NoError(t, err)
			assert.Equal(t, tt.want.uid, uuid.String())
		})
	}
}
