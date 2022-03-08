package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postURLApiHandler(t *testing.T) {
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
		baseURL string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/api/shorten",
			longURL: "https://github.com/serjyuriev/",
			baseURL: "http://localhost:8080",
			want: want{
				statusCode:  201,
				contentType: "application/json",
				generateURL: true,
				urlRegex:    regexp.MustCompile("http://localhost:8080/[a-z]{6}"),
			},
		},
		{
			name:    "empty body",
			request: "http://localhost:8080/api/shorten",
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
			request: "http://localhost:8080/api/shorten",
			longURL: "wow",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				generateURL: false,
				response:    "bad request\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ShortURLHost = tt.baseURL
			Store, err = storage.NewStore("")
			require.NoError(t, err)
			reqBody := postShortenRequest{
				URL: tt.longURL,
			}
			reqBz, err := json.Marshal(reqBody)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(reqBz))
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostURLApiHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if tt.want.generateURL {
				var res postShortenResponse
				err := json.NewDecoder(result.Body).Decode(&res)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.Regexp(t, tt.want.urlRegex, res.Result)
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
		baseURL string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/",
			longURL: "https://github.com/serjyuriev/",
			baseURL: "http://localhost:8080",
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
				response:    "bad request\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ShortURLHost = tt.baseURL
			Store, err = storage.NewStore("")
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.longURL))
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostURLHandler)
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
			var err error
			Store, err = storage.NewStore("")
			require.NoError(t, err)
			uid := uuid.New().String()
			Store.InsertNewURLPair(uid, "abcdef", "https://github.com/serjyuriev/")
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetURLHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
