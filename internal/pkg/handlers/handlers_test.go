package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/serjyuriev/shortener/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postURLHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		generateURL bool
		urlRegex    *regexp.Regexp
		response    string
	}
	tests := []struct {
		name    string
		request string
		longURL string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/",
			longURL: "https://github.com/serjyuriev/",
			want: want{
				statusCode:  201,
				contentType: "text/plain",
				generateURL: true,
				urlRegex:    regexp.MustCompile("http://localhost:8080/[a-z]{6}"),
			},
		},
		{
			name:    "empty body",
			request: "http://localhost:8080/",
			longURL: "",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				generateURL: false,
				response:    "Body cannot be empty.\n",
			},
		},
		{
			name:    "not URL in body",
			request: "http://localhost:8080/",
			longURL: "wow",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				generateURL: false,
				response:    "parse \"wow\": invalid URI for request\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, _ = storage.NewStore(context.Background())
			cfg = serverConfig{
				host: "localhost",
				port: 8080,
			}
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.longURL))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(postURLHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if tt.want.generateURL {
				response, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.Regexp(t, tt.want.urlRegex, string(response))
			} else {
				response, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.Equal(t, tt.want.response, string(response))
			}

		})
	}
}

func Test_getURLHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/abcdef",
			want: want{
				statusCode: 307,
				location:   "https://github.com/serjyuriev/",
			},
		},
		{
			name:    "no path",
			request: "http://localhost:8080/",
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
		{
			name:    "wrong short URL",
			request: "http://localhost:8080/fedcba",
			want: want{
				statusCode: 400,
				location:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, _ = storage.NewStore(context.Background())
			store.InsertNewURLPair(context.Background(), storage.ShortPath("abcdef"), storage.LongURL("https://github.com/serjyuriev/"))
			cfg = serverConfig{
				host: "localhost",
				port: 8080,
			}
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(getURLHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
