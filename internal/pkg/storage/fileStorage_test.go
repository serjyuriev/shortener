package storage

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	arrayStore = "array"
	mapStore   = "map"
)

func Test_fileStore_DeleteManyURLs(t *testing.T) {
	tests := []struct {
		name      string
		storeType string
	}{
		{
			name:      "delete many (map)",
			storeType: mapStore,
		},
		{
			name:      "delete many (array)",
			storeType: arrayStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s   Store
				err error
			)
			switch tt.storeType {
			case mapStore:
				s, err = NewFileStore("")
			case arrayStore:
				s, err = NewFileArrayStore("")
			}
			require.NoError(t, err)

			err = s.DeleteManyURLs(context.Background(), uuid.Nil, nil)
			require.ErrorIs(t, err, ErrNotImplementedYet)
		})
	}
}

func Test_fileStore_FindByOriginalURL(t *testing.T) {
	type want struct {
		hasError bool
		shortURL string
	}
	tests := []struct {
		name        string
		storeType   string
		originalURL string
		want        want
	}{
		{
			name:        "existing original URL is provided (map)",
			storeType:   mapStore,
			originalURL: "https://yandex.ru",
			want: want{
				hasError: false,
				shortURL: "lkasdj",
			},
		},
		{
			name:        "non-existing original URL is provided (map)",
			storeType:   mapStore,
			originalURL: "https://duckduckgo.com",
			want: want{
				hasError: true,
				shortURL: "",
			},
		},
		{
			name:        "no original URL is provided (map)",
			storeType:   mapStore,
			originalURL: "",
			want: want{
				hasError: true,
				shortURL: "",
			},
		},
		{
			name:        "existing original URL is provided (array)",
			storeType:   arrayStore,
			originalURL: "https://yandex.ru",
			want: want{
				hasError: false,
				shortURL: "lkasdj",
			},
		},
		{
			name:        "non-existing original URL is provided (array)",
			storeType:   arrayStore,
			originalURL: "https://duckduckgo.com",
			want: want{
				hasError: true,
				shortURL: "",
			},
		},
		{
			name:        "no original URL is provided (array)",
			storeType:   arrayStore,
			originalURL: "",
			want: want{
				hasError: true,
				shortURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s   Store
				err error
			)
			switch tt.storeType {
			case mapStore:
				s, err = NewFileStore("")
			case arrayStore:
				s, err = NewFileArrayStore("")
			}
			require.NoError(t, err)

			uid := uuid.New()
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

			err = s.InsertManyURLs(context.Background(), uid, urls)
			require.NoError(t, err)

			url, err := s.FindByOriginalURL(context.Background(), tt.originalURL)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.shortURL, url)
		})
	}
}

func Test_fileStore_FindOriginalURL(t *testing.T) {
	type want struct {
		hasError    bool
		originalURL string
	}
	tests := []struct {
		name      string
		storeType string
		shortURL  string
		want      want
	}{
		{
			name:      "existing short URL is provided (map)",
			storeType: mapStore,
			shortURL:  "sdfkbj",
			want: want{
				hasError:    false,
				originalURL: "https://discord.com",
			},
		},
		{
			name:      "non-existing short URL is provided (map)",
			storeType: mapStore,
			shortURL:  "ashbxc",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
		{
			name:      "no short URL is provided (map)",
			storeType: mapStore,
			shortURL:  "",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
		{
			name:      "existing short URL is provided (array)",
			storeType: arrayStore,
			shortURL:  "sdfkbj",
			want: want{
				hasError:    false,
				originalURL: "https://discord.com",
			},
		},
		{
			name:      "non-existing short URL is provided (array)",
			storeType: arrayStore,
			shortURL:  "ashbxc",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
		{
			name:      "no short URL is provided (array)",
			storeType: arrayStore,
			shortURL:  "",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s   Store
				err error
			)
			switch tt.storeType {
			case mapStore:
				s, err = NewFileStore("")
			case arrayStore:
				s, err = NewFileArrayStore("")
			}
			require.NoError(t, err)

			uid := uuid.New()
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

			err = s.InsertManyURLs(context.Background(), uid, urls)
			require.NoError(t, err)

			url, err := s.FindOriginalURL(context.Background(), tt.shortURL)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.originalURL, url)
		})
	}
}

func Test_fileStore_FindURLsByUser(t *testing.T) {
	type want struct {
		hasError bool
		err      error
		urls     map[string]string
	}
	tests := []struct {
		name      string
		storeType string
		userID    uuid.UUID
		userURLs  map[string]string
		otherURLs map[string]string
		want      want
	}{
		{
			name:      "correct UUID of existing user is provided (map)",
			storeType: mapStore,
			userID:    uuid.New(),
			userURLs: map[string]string{
				"abcdef": "https://github.com/serjyuriev",
				"fedcba": "https://gitlab.com/servady",
				"zkljns": "https://twitch.tv",
			},
			otherURLs: map[string]string{
				"lkasdj": "https://yandex.ru",
				"aslkqs": "https://google.com",
				"cpsoks": "https://vk.com",
				"qwrkml": "https://habr.com",
				"sdfkbj": "https://discord.com",
				"qlwknf": "https://gmail.com",
				"qkwnmd": "https://vscode.dev",
			},
			want: want{
				hasError: false,
				err:      nil,
				urls: map[string]string{
					"abcdef": "https://github.com/serjyuriev",
					"fedcba": "https://gitlab.com/servady",
					"zkljns": "https://twitch.tv",
				},
			},
		},
		{
			name:      "correct UUID of non-existing user is provided (map)",
			storeType: mapStore,
			userID:    uuid.New(),
			userURLs:  map[string]string{},
			otherURLs: map[string]string{
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
			},
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				urls:     nil,
			},
		},
		{
			name:      "no UUID is provided (map)",
			storeType: mapStore,
			userID:    uuid.Nil,
			userURLs:  map[string]string{},
			otherURLs: map[string]string{
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
			},
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				urls:     nil,
			},
		},
		{
			name:      "correct UUID of existing user is provided (array)",
			storeType: arrayStore,
			userID:    uuid.New(),
			userURLs: map[string]string{
				"abcdef": "https://github.com/serjyuriev",
				"fedcba": "https://gitlab.com/servady",
				"zkljns": "https://twitch.tv",
			},
			otherURLs: map[string]string{
				"lkasdj": "https://yandex.ru",
				"aslkqs": "https://google.com",
				"cpsoks": "https://vk.com",
				"qwrkml": "https://habr.com",
				"sdfkbj": "https://discord.com",
				"qlwknf": "https://gmail.com",
				"qkwnmd": "https://vscode.dev",
			},
			want: want{
				hasError: false,
				err:      nil,
				urls: map[string]string{
					"abcdef": "https://github.com/serjyuriev",
					"fedcba": "https://gitlab.com/servady",
					"zkljns": "https://twitch.tv",
				},
			},
		},
		{
			name:      "correct UUID of non-existing user is provided (array)",
			storeType: arrayStore,
			userID:    uuid.New(),
			userURLs:  map[string]string{},
			otherURLs: map[string]string{
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
			},
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				urls:     nil,
			},
		},
		{
			name:      "no UUID is provided (array)",
			storeType: arrayStore,
			userID:    uuid.Nil,
			userURLs:  map[string]string{},
			otherURLs: map[string]string{
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
			},
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				urls:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s   Store
				err error
			)
			switch tt.storeType {
			case mapStore:
				s, err = NewFileStore("")
			case arrayStore:
				s, err = NewFileArrayStore("")
			}
			require.NoError(t, err)

			uid := uuid.New()
			err = s.InsertManyURLs(context.Background(), uid, tt.otherURLs)
			require.NoError(t, err)
			err = s.InsertManyURLs(context.Background(), tt.userID, tt.userURLs)
			require.NoError(t, err)

			urls, err := s.FindURLsByUser(context.Background(), tt.userID)
			if tt.want.hasError {
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.urls, urls)
		})
	}
}

func Test_fileStore_InsertManyURLs(t *testing.T) {
	type want struct {
		hasError bool
		err      error
		length   int
	}
	tests := []struct {
		name           string
		storeType      string
		userID         uuid.UUID
		urls           map[string]string
		useFileStorage bool
		storagePath    string
		hasWriteError  bool
		want           want
	}{
		{
			name:      "correct data is provided (map)",
			storeType: mapStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: false,
			storagePath:    "",
			hasWriteError:  false,
			want: want{
				hasError: false,
				err:      nil,
				length:   10,
			},
		},
		{
			name:      "correct data is provided, file storage (map)",
			storeType: mapStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  false,
			want: want{
				hasError: false,
				err:      nil,
				length:   10,
			},
		},
		{
			name:      "correct data is provided, file storage write error (map)",
			storeType: mapStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  true,
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				length:   0,
			},
		},
		{
			name:      "correct data is provided (array)",
			storeType: arrayStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: false,
			storagePath:    "",
			hasWriteError:  false,
			want: want{
				hasError: false,
				err:      nil,
				length:   10,
			},
		},
		{
			name:      "correct data is provided, file storage (array)",
			storeType: arrayStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  false,
			want: want{
				hasError: false,
				err:      nil,
				length:   10,
			},
		},
		{
			name:      "correct data is provided, file storage write error (array)",
			storeType: arrayStore,
			userID:    uuid.New(),
			urls: map[string]string{
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
			},
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  true,
			want: want{
				hasError: true,
				err:      ErrNoURLWasFound,
				length:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				store Store
				err   error
			)
			if tt.useFileStorage {
				if file, err := os.Open(tt.storagePath); err == nil {
					file.Close()
					err = os.Remove(tt.storagePath)
					require.NoError(t, err)
				}
			}

			switch tt.storeType {
			case mapStore:
				var s *fileStore
				if tt.useFileStorage {
					s = &fileStore{
						fileStoragePath: tt.storagePath,
						useFileStorage:  tt.useFileStorage,
					}
				} else {
					s = &fileStore{
						fileStoragePath: "",
						useFileStorage:  false,
					}
				}
				s.URLs = make(map[string]link)
				if tt.useFileStorage && tt.hasWriteError {
					s.fileStoragePath = ""
				}
				store = s
			case arrayStore:
				var s *fileArrayStore
				if tt.useFileStorage {
					s = &fileArrayStore{
						fileStoragePath: tt.storagePath,
						useFileStorage:  tt.useFileStorage,
					}
				} else {
					s = &fileArrayStore{
						fileStoragePath: "",
						useFileStorage:  false,
					}
				}
				s.URLs = make([]arrayLink, 0)
				if tt.useFileStorage && tt.hasWriteError {
					s.fileStoragePath = ""
				}
				store = s
			}

			err = store.InsertManyURLs(
				context.Background(),
				tt.userID,
				tt.urls,
			)
			if tt.want.hasError {
				assert.Error(t, err)
				_, err = store.FindURLsByUser(context.Background(), tt.userID)
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				assert.NoError(t, err)
				urls, err := store.FindURLsByUser(context.Background(), tt.userID)
				require.NoError(t, err)
				assert.Equal(t, tt.want.length, len(urls))
			}
		})
	}
}

func Test_fileStore_InsertNewURLPair(t *testing.T) {
	type want struct {
		hasError   bool
		err        error
		orignalURL string
	}
	tests := []struct {
		name           string
		storeType      string
		userID         uuid.UUID
		shortPath      string
		originalURL    string
		useFileStorage bool
		storagePath    string
		hasWriteError  bool
		want           want
	}{
		{
			name:           "correct data is provided (map)",
			storeType:      mapStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: false,
			storagePath:    "",
			hasWriteError:  false,
			want: want{
				hasError:   false,
				err:        nil,
				orignalURL: "https://duckduckgo.com",
			},
		},
		{
			name:           "correct data is provided, file storage (map)",
			storeType:      mapStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  false,
			want: want{
				hasError:   false,
				err:        nil,
				orignalURL: "https://duckduckgo.com",
			},
		},
		{
			name:           "correct data is provided, file storage write error (map)",
			storeType:      mapStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  true,
			want: want{
				hasError:   true,
				err:        ErrNoURLWasFound,
				orignalURL: "",
			},
		},
		{
			name:           "correct data is provided (array)",
			storeType:      arrayStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: false,
			storagePath:    "",
			hasWriteError:  false,
			want: want{
				hasError:   false,
				err:        nil,
				orignalURL: "https://duckduckgo.com",
			},
		},
		{
			name:           "correct data is provided, file storage (array)",
			storeType:      arrayStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  false,
			want: want{
				hasError:   false,
				err:        nil,
				orignalURL: "https://duckduckgo.com",
			},
		},
		{
			name:           "correct data is provided, file storage write error (array)",
			storeType:      arrayStore,
			userID:         uuid.New(),
			shortPath:      "abcdef",
			originalURL:    "https://duckduckgo.com",
			useFileStorage: true,
			storagePath:    "shorty.json",
			hasWriteError:  true,
			want: want{
				hasError:   true,
				err:        ErrNoURLWasFound,
				orignalURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				store Store
				err   error
			)
			if tt.useFileStorage {
				if file, err := os.Open(tt.storagePath); err == nil {
					file.Close()
					err = os.Remove(tt.storagePath)
					require.NoError(t, err)
				}
			}

			switch tt.storeType {
			case mapStore:
				var s *fileStore
				if tt.useFileStorage {
					s = &fileStore{
						fileStoragePath: tt.storagePath,
						useFileStorage:  tt.useFileStorage,
					}
				} else {
					s = &fileStore{
						fileStoragePath: "",
						useFileStorage:  false,
					}
				}
				s.URLs = make(map[string]link)
				if tt.useFileStorage && tt.hasWriteError {
					s.fileStoragePath = ""
				}
				store = s
			case arrayStore:
				var s *fileArrayStore
				if tt.useFileStorage {
					s = &fileArrayStore{
						fileStoragePath: tt.storagePath,
						useFileStorage:  tt.useFileStorage,
					}
				} else {
					s = &fileArrayStore{
						fileStoragePath: "",
						useFileStorage:  false,
					}
				}
				s.URLs = make([]arrayLink, 0)
				if tt.useFileStorage && tt.hasWriteError {
					s.fileStoragePath = ""
				}
				store = s
			}

			err = store.InsertNewURLPair(
				context.Background(),
				tt.userID,
				tt.shortPath,
				tt.originalURL,
			)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orig, err := store.FindOriginalURL(context.Background(), tt.shortPath)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.orignalURL, orig)
		})
	}
}

func Test_fileStore_Ping(t *testing.T) {
	type want struct {
		hasError bool
	}
	tests := []struct {
		name      string
		storeType string
		want      want
	}{
		{
			name:      "just ping (map)",
			storeType: mapStore,
			want: want{
				hasError: false,
			},
		},
		{
			name:      "just ping (array)",
			storeType: arrayStore,
			want: want{
				hasError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s   Store
				err error
			)
			switch tt.storeType {
			case mapStore:
				s, err = NewFileStore("")
			case arrayStore:
				s, err = NewFileArrayStore("")
			}
			require.NoError(t, err)
			err = s.Ping(context.Background())
			require.NoError(t, err)
		})
	}
}

func Test_fileStore_loadDataFromFile(t *testing.T) {
	type want struct {
		hasError    bool
		originalURL string
	}
	tests := []struct {
		name        string
		storeType   string
		json        string
		storagePath string
		want        want
	}{
		{
			name:        "data provided (map)",
			storeType:   mapStore,
			json:        "{\"abcdef\":{\"Original\":\"https://github.com/serjyuriev\", \"User\":\"8ebc62e1-63d2-4cc1-b8cf-20cdcc797f3c\"}}",
			storagePath: "shorty.json",
			want: want{
				hasError:    false,
				originalURL: "https://github.com/serjyuriev",
			},
		},
		{
			name:        "empty file (map)",
			storeType:   mapStore,
			json:        "",
			storagePath: "shorty.json",
			want: want{
				hasError:    false,
				originalURL: "",
			},
		},
		{
			name:        "not json (map)",
			storeType:   mapStore,
			json:        "]fa s.df ]2qe[f. ][,f 1,",
			storagePath: "shorty.json",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
		{
			name:        "data provided (array)",
			storeType:   arrayStore,
			json:        "[{\"Original\":\"https://github.com/serjyuriev\",\"User\":\"8ebc62e1-63d2-4cc1-b8cf-20cdcc797f3c\",\"Shortened\":\"abcdef\"}]",
			storagePath: "shorty.json",
			want: want{
				hasError:    false,
				originalURL: "https://github.com/serjyuriev",
			},
		},
		{
			name:        "empty file (array)",
			storeType:   arrayStore,
			json:        "",
			storagePath: "shorty.json",
			want: want{
				hasError:    false,
				originalURL: "",
			},
		},
		{
			name:        "not json (array)",
			storeType:   arrayStore,
			json:        "]fa s.df ]2qe[f. ][,f 1,",
			storagePath: "shorty.json",
			want: want{
				hasError:    true,
				originalURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.storeType == mapStore {
				s := &fileStore{
					fileStoragePath: tt.storagePath,
					useFileStorage:  true,
				}

				file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
				require.NoError(t, err)
				_, err = file.WriteString(tt.json)
				require.NoError(t, err)
				err = file.Close()
				require.NoError(t, err)

				err = s.loadDataFromFile()
				if tt.want.hasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if tt.json != "" {
						assert.Equal(t, tt.want.originalURL, s.URLs["abcdef"].Original)
					}
				}
			} else {
				s := &fileArrayStore{
					fileStoragePath: tt.storagePath,
					useFileStorage:  true,
				}

				file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
				require.NoError(t, err)
				_, err = file.WriteString(tt.json)
				require.NoError(t, err)
				err = file.Close()
				require.NoError(t, err)

				err = s.loadDataFromFile()
				if tt.want.hasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if tt.json != "" {
						assert.Equal(t, tt.want.originalURL, s.URLs[0].Original)
					}
				}
			}

			err := os.Remove(tt.storagePath)
			require.NoError(t, err)
		})
	}
}

func Test_fileStore_writeDataToFile(t *testing.T) {
	type want struct {
		hasError bool
	}
	tests := []struct {
		name        string
		userID      uuid.UUID
		urls        map[string]string
		storagePath string
		want        want
	}{
		{
			name:   "path is provided, urls initialized",
			userID: uuid.New(),
			urls: map[string]string{
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
			},
			storagePath: "shorty.json",
			want: want{
				hasError: false,
			},
		},
		{
			name:   "no path is provided, urls are initialized",
			userID: uuid.New(),
			urls: map[string]string{
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
			},
			storagePath: "",
			want: want{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fileStore{
				fileStoragePath: tt.storagePath,
				useFileStorage:  true,
			}

			s.URLs = make(map[string]link)
			for k, v := range tt.urls {
				s.URLs[k] = link{
					Original: v,
					User:     tt.userID,
				}
			}

			err := s.writeDataToFile()
			if tt.want.hasError {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			s2 := &fileStore{
				fileStoragePath: tt.storagePath,
				useFileStorage:  true,
			}
			err = s2.loadDataFromFile()
			require.NoError(t, err)
			assert.Equal(t, s.URLs, s2.URLs)

			err = os.Remove(s.fileStoragePath)
			require.NoError(t, err)
		})
	}
}
