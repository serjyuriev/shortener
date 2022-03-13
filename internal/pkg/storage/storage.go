package storage

import (
	"context"
	"errors"
)

var (
	ErrNotImplementedYet = errors.New("method not implemented yet")
	ErrNoURLWasFound     = errors.New("no URL was found")
)

type Store interface {
	FindLongURL(ctx context.Context, shortPath string) (string, error)
	FindURLsByUser(ctx context.Context, userID string) (map[string]string, error)
	InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
}
