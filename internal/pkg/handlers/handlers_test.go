package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/serjyuriev/shortener/internal/pkg/service"
)

func ExampleHandlers_DeleteURLsHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}
	if err = h.svc.InsertManyURLs(context.Background(), uid, urls); err != nil {
		fmt.Printf("unable to insert urls: %v\n", err)
		return
	}

	toDelete := []string{
		"qwrkml",
		"fedcba",
	}
	reqBody, err := json.Marshal(toDelete)
	if err != nil {
		fmt.Printf("unable to marshal request body: %v\n", err)
		return
	}
	request, err := http.NewRequest(
		http.MethodDelete,
		"http://localhost:8080/api/user/urls",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			contextKeyUID,
			uid,
		),
	)

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.DeleteURLsHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	fmt.Printf("Code: %d\n", result.StatusCode)

	// Output:
	// Code: 202
}

func ExampleHandlers_GetURLHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}
	if err = h.svc.InsertManyURLs(context.Background(), uid, urls); err != nil {
		fmt.Printf("unable to insert urls: %v\n", err)
		return
	}

	request, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/zkljns",
		nil,
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.GetURLHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	fmt.Printf("Code: %d\n", result.StatusCode)
	fmt.Printf("Original URL: %s\n", result.Header.Get("Location"))

	// Output:
	// Code: 307
	// Original URL: https://twitch.tv
}

func ExampleHandlers_GetUserURLsAPIHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
	}

	uid2 := uuid.New().String()
	urls2 := map[string]string{
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}
	if err = h.svc.InsertManyURLs(context.Background(), uid, urls); err != nil {
		fmt.Printf("unable to insert urls: %v\n", err)
		return
	}
	if err = h.svc.InsertManyURLs(context.Background(), uid2, urls2); err != nil {
		fmt.Printf("unable to insert urls: %v\n", err)
		return
	}

	request, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/api/user/urls",
		nil,
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			contextKeyUID,
			uid,
		),
	)

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.GetUserURLsAPIHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	var firstUserURLs []userURLs
	err = json.NewDecoder(result.Body).Decode(&firstUserURLs)

	fmt.Printf("Code: %d\n", result.StatusCode)
	fmt.Printf("URLs found: %d\n", len(firstUserURLs))

	// Output:
	// Code: 200
	// URLs found: 5
}

func ExampleHandlers_PingHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	request, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/ping",
		nil,
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.PingHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("unable to read response body: %v\n", err)
		return
	}

	fmt.Printf("Code: %d\n", result.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	// Output:
	// Code: 200
	// Response: ok
}

func ExampleHandlers_PostBatchHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	reqArray := []postBatchSingleRequest{
		{
			CorrelationID: "12345",
			OriginalURL:   "https://github.com/serjyuriev",
		},
		{
			CorrelationID: "12345",
			OriginalURL:   "https://gitlab.com/servady",
		},
		{
			CorrelationID: "12345",
			OriginalURL:   "https://google.com",
		},
		{
			CorrelationID: "12345",
			OriginalURL:   "https://twitch.tv",
		},
		{
			CorrelationID: "12345",
			OriginalURL:   "https://habr.com",
		},
	}
	reqBody, err := json.Marshal(reqArray)
	if err != nil {
		fmt.Printf("unable to marshal request body: %v\n", err)
		return
	}
	request, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten/batch",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			contextKeyUID,
			uid,
		),
	)

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.PostBatchHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	var response []postBatchSingleResponse
	if err = json.NewDecoder(result.Body).Decode(&response); err != nil {
		fmt.Printf("unable to decode response body: %v\n", err)
		return
	}

	fmt.Printf("Code: %d\n", result.StatusCode)
	fmt.Printf("Added URLs: %d\n", len(response))

	// Output:
	// Code: 201
	// Added URLs: 5
}

func ExampleHandlers_PostURLApiHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	req := postShortenRequest{URL: "https://twitch.tv"}
	reqBody, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("unable to marshal request body: %v\n", err)
		return
	}
	request, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			contextKeyUID,
			uid,
		),
	)

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.PostURLApiHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	fmt.Printf("Code: %d\n", result.StatusCode)

	// Output:
	// Code: 201
}

func ExampleHandlers_PostURLHandler() {
	svc, err := service.NewService()
	if err != nil {
		fmt.Printf("unable to initiazlize service: %v\n", err)
		return
	}
	h := &Handlers{svc: svc}

	uid := uuid.New().String()
	url := []byte("https://github.com/serjyuriev")
	request, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		bytes.NewBuffer(url),
	)
	if err != nil {
		fmt.Printf("unable to initialize http request: %v\n", err)
		return
	}
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			contextKeyUID,
			uid,
		),
	)

	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.PostURLHandler)
	hf.ServeHTTP(w, request)
	result := w.Result()
	defer result.Body.Close()

	fmt.Printf("Code: %d\n", result.StatusCode)

	// Output:
	// Code: 201
}

func TestPostBatchHandler(t *testing.T) {
	type want struct {
		statusCode    int
		contentType   string
		generateURL   bool
		urlRegex      *regexp.Regexp
		correlationID string
	}
	tests := []struct {
		name          string
		request       string
		longURL       []string
		correlationID string
		baseURL       string
		want          want
	}{
		{
			name:    "positive test #1",
			request: "http://localhost:8080/api/shorten/batch",
			longURL: []string{
				"https://vk.com",
				"https://youtube.com",
				"https://google.com",
				"https://gitlab.com",
				"https://github.com",
			},
			correlationID: "12345",
			baseURL:       "http://localhost:8080",
			want: want{
				statusCode:    201,
				contentType:   "application/json",
				generateURL:   true,
				urlRegex:      regexp.MustCompile("http://localhost:8080/[a-z]{6}"),
				correlationID: "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := service.NewService()
			require.NoError(t, err)
			h := &Handlers{
				baseURL: tt.baseURL,
				svc:     svc,
			}
			reqBody := make([]postBatchSingleRequest, 0)
			for _, url := range tt.longURL {
				reqBody = append(
					reqBody,
					postBatchSingleRequest{
						CorrelationID: tt.correlationID,
						OriginalURL:   url,
					},
				)
			}
			reqBz, err := json.Marshal(reqBody)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(reqBz))
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
			w := httptest.NewRecorder()
			hf := http.HandlerFunc(h.PostBatchHandler)
			hf.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if tt.want.generateURL {
				var res []postBatchSingleResponse
				err := json.NewDecoder(result.Body).Decode(&res)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)
				assert.Regexp(t, tt.want.urlRegex, res[0].ShortURL)
				assert.Equal(t, tt.want.correlationID, res[0].CorrelationID)
			}
		})
	}
}

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
			svc, err := service.NewService()
			require.NoError(t, err)
			h := &Handlers{
				baseURL: tt.baseURL,
				svc:     svc,
			}
			reqBody := postShortenRequest{
				URL: tt.longURL,
			}
			reqBz, err := json.Marshal(reqBody)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(reqBz))
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
			w := httptest.NewRecorder()
			// h := http.HandlerFunc(PostURLApiHandler)
			hf := http.HandlerFunc(h.PostURLApiHandler)
			hf.ServeHTTP(w, request)
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
			svc, err := service.NewService()
			require.NoError(t, err)
			h := &Handlers{
				baseURL: tt.baseURL,
				svc:     svc,
			}
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.longURL))
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
			w := httptest.NewRecorder()
			hf := http.HandlerFunc(h.PostURLHandler)
			hf.ServeHTTP(w, request)
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
			svc, err := service.NewService()
			require.NoError(t, err)
			h := &Handlers{
				svc: svc,
			}
			uid := uuid.New().String()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			h.svc.InsertNewURLPair(ctx, uid, "abcdef", "https://github.com/serjyuriev/")
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			hf := http.HandlerFunc(h.GetURLHandler)
			hf.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func TestGetUserURLsAPIHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		response    []userURLs
	}
	tests := []struct {
		name    string
		baseURL string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			baseURL: "http://localhost:8080",
			request: "http://localhost:8080/api/user/urls",
			want: want{
				statusCode:  200,
				contentType: "application/json",
				response: []userURLs{
					{
						ShortURL:    "http://localhost:8080/lizuyl",
						OriginalURL: "https://gitlab.com",
					},
					{
						ShortURL:    "http://localhost:8080/ppgcni",
						OriginalURL: "https://vk.com",
					},
					{
						ShortURL:    "http://localhost:8080/ugkqzj",
						OriginalURL: "https://github.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := service.NewService()
			require.NoError(t, err)
			h := &Handlers{
				svc:     svc,
				baseURL: tt.baseURL,
			}
			uid := uuid.New().String()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			h.svc.InsertNewURLPair(ctx, uid, "lizuyl", "https://gitlab.com")
			h.svc.InsertNewURLPair(ctx, uid, "ppgcni", "https://vk.com")
			h.svc.InsertNewURLPair(ctx, uid, "ugkqzj", "https://github.com")
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uid))
			w := httptest.NewRecorder()
			hf := http.HandlerFunc(h.GetUserURLsAPIHandler)
			hf.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			var urls []userURLs
			err = json.NewDecoder(result.Body).Decode(&urls)
			require.NoError(t, err)
			for _, url := range urls {
				assert.Contains(t, tt.want.response, url)
			}
		})
	}
}

func BenchmarkGetURLHandler(b *testing.B) {
	svc, err := service.NewService()
	if err != nil {
		b.Errorf("unable to initiazlize service: %v\n", err)
	}
	h := &Handlers{
		svc: svc,
	}
	uid := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	h.svc.InsertNewURLPair(ctx, uid, "abcdef", "https://vk.com/groups")
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/abcdef", nil)
	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.GetURLHandler)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hf.ServeHTTP(w, request)
	}
}

func BenchmarkPostURLHandler(b *testing.B) {
	svc, _ := service.NewService()
	h := &Handlers{
		baseURL: "http://localhost:8080",
		svc:     svc,
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		strings.NewReader("https://gitlab.com/servady"),
	)
	request = request.WithContext(context.WithValue(request.Context(), contextKeyUID, uuid.New().String()))
	w := httptest.NewRecorder()
	hf := http.HandlerFunc(h.PostURLHandler)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hf.ServeHTTP(w, request)
	}
}
