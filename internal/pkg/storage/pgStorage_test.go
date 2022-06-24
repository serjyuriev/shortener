package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteManyURLs(t *testing.T) {
	userID := uuid.New()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://shorty:sh0rt4@localhost:5432/shortenertest"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Logf("unable to connect to postgre: %v\n", err)
		t.FailNow()
	}

	s := &pgStore{
		cs: dsn,
		db: db,
	}

	_, err = s.db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
			short_id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			added_by_user TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`,
	)
	if err != nil {
		t.Logf("unable to create table: %v\n", err)
		t.FailNow()
	}

	err = s.InsertManyURLs(
		context.Background(),
		userID,
		map[string]string{
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
	)
	if err != nil {
		t.Logf("unable to insert values: %v\n", err)
		t.FailNow()
	}

	urls := []string{
		"abcdef",
		"fedcba",
		"lkasdj",
		"aslkqs",
	}
	wantLength := 6

	err = s.DeleteManyURLs(context.Background(), userID, urls)
	assert.NoError(t, err)

	orig, err := s.FindURLsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, wantLength, len(orig))

	_, err = s.db.Exec("DROP INDEX IF EXISTS original_url_idx; DROP TABLE IF EXISTS urls;")
	if err != nil {
		t.Logf("unable to drop table: %v\n", err)
	}
}

func TestFindByOriginalURL(t *testing.T) {
	userID := uuid.New()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://shorty:sh0rt4@localhost:5432/shortenertest"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Logf("unable to connect to postgre: %v\n", err)
		t.FailNow()
	}

	s := &pgStore{
		cs: dsn,
		db: db,
	}

	_, err = s.db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
			short_id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			added_by_user TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`,
	)
	if err != nil {
		t.Logf("unable to create table: %v\n", err)
		t.FailNow()
	}

	err = s.InsertManyURLs(
		context.Background(),
		userID,
		map[string]string{
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
	)
	if err != nil {
		t.Logf("unable to insert values: %v\n", err)
		t.FailNow()
	}

	type want struct {
		hasError bool
	}
	tests := []struct {
		name     string
		userID   uuid.UUID
		original string
		delete   bool
		short    string
		want     want
	}{
		{
			name:     "existing original URL",
			userID:   userID,
			original: "https://discord.com",
			delete:   false,
			short:    "sdfkbj",
			want: want{
				hasError: false,
			},
		},
		{
			name:     "deleted original URL",
			userID:   userID,
			original: "https://discord.com",
			delete:   true,
			short:    "sdfkbj",
			want: want{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.delete {
				err = s.DeleteManyURLs(context.Background(), userID, []string{tt.short})
				require.NoError(t, err)
			}

			short, err := s.FindByOriginalURL(context.Background(), tt.original)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.short, short)
			}
		})
	}

	_, err = s.db.Exec("DROP INDEX IF EXISTS original_url_idx; DROP TABLE IF EXISTS urls;")
	if err != nil {
		t.Logf("unable to drop table: %v\n", err)
	}
}

func TestInsertManyURLs(t *testing.T) {
	userID := uuid.New()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://shorty:sh0rt4@localhost:5432/shortenertest"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Logf("unable to connect to postgre: %v\n", err)
		t.FailNow()
	}

	s := &pgStore{
		cs: dsn,
		db: db,
	}

	_, err = s.db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
			short_id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			added_by_user TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`,
	)
	if err != nil {
		t.Logf("unable to create table: %v\n", err)
		t.FailNow()
	}

	type want struct {
		hasError bool
		length   int
	}
	tests := []struct {
		name   string
		userID uuid.UUID
		urls   map[string]string
		want   want
	}{
		{
			name:   "new original URL",
			userID: userID,
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
			want: want{
				hasError: false,
				length:   10,
			},
		},
		{
			name:   "already existing short path",
			userID: userID,
			urls: map[string]string{
				"abcdef": "https://github.com/serjyuriev",
			},
			want: want{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = s.InsertManyURLs(context.Background(), tt.userID, tt.urls)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				orig, err := s.FindURLsByUser(context.Background(), tt.userID)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.length, len(orig))
			}
		})
	}

	_, err = s.db.Exec("DROP INDEX IF EXISTS original_url_idx; DROP TABLE IF EXISTS urls;")
	if err != nil {
		t.Logf("unable to drop table: %v\n", err)
	}
}

func TestInsertNewURLPair(t *testing.T) {
	userID := uuid.New()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://shorty:sh0rt4@localhost:5432/shortenertest"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Logf("unable to connect to postgre: %v\n", err)
		t.FailNow()
	}

	s := &pgStore{
		cs: dsn,
		db: db,
	}

	_, err = s.db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
			short_id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			added_by_user TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`,
	)
	if err != nil {
		t.Logf("unable to create table: %v\n", err)
		t.FailNow()
	}

	type want struct {
		hasError bool
	}
	tests := []struct {
		name        string
		userID      uuid.UUID
		shortPath   string
		originalURL string
		want        want
	}{
		{
			name:        "new original URL",
			userID:      userID,
			shortPath:   "abcdef",
			originalURL: "https://duckduckgo.com",
			want: want{
				hasError: false,
			},
		},
		{
			name:        "already existing short path",
			userID:      userID,
			shortPath:   "abcdef",
			originalURL: "https://twitch.tv",
			want: want{
				hasError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = s.InsertNewURLPair(
				context.Background(),
				tt.userID,
				tt.shortPath,
				tt.originalURL,
			)
			if tt.want.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				orig, err := s.FindOriginalURL(context.Background(), tt.shortPath)
				assert.NoError(t, err)
				assert.Equal(t, tt.originalURL, orig)
			}
		})
	}

	_, err = s.db.Exec("DROP INDEX IF EXISTS original_url_idx; DROP TABLE IF EXISTS urls;")
	if err != nil {
		t.Logf("unable to drop table: %v\n", err)
	}
}

func TestPing(t *testing.T) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://shorty:sh0rt4@localhost:5432/shortenertest"
	}
	s, err := NewPgStore(dsn)
	require.NoError(t, err)
	err = s.Ping(context.Background())
	assert.NoError(t, err)
}
