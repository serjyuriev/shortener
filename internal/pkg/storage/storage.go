package storage

import (
	"context"
	"errors"
)

var (
	ErrNotImplementedYet    = errors.New("method not implemented yet")
	ErrNoURLWasFound        = errors.New("no URL was found")
	ErrNotUniqueOriginalURL = errors.New("original URL already presented")
)

type Store interface {
	FindByLongURL(ctx context.Context, long string) (string, error)
	FindLongURL(ctx context.Context, shortPath string) (string, error)
	FindURLsByUser(ctx context.Context, userID string) (map[string]string, error)
	InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error
	InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
}
