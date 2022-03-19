package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotImplementedYet    = errors.New("method not implemented yet")
	ErrNoURLWasFound        = errors.New("no URL was found")
	ErrNotUniqueOriginalURL = errors.New("original URL already presented")
)

type Store interface {
	FindByOriginalURL(ctx context.Context, originalURL string) (string, error)
	FindOriginalURL(ctx context.Context, shortPath string) (string, error)
	FindURLsByUser(ctx context.Context, userID uuid.UUID) (map[string]string, error)
	InsertManyURLs(ctx context.Context, userID uuid.UUID, urls map[string]string) error
	InsertNewURLPair(ctx context.Context, userID uuid.UUID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
}
