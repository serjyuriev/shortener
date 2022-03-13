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
	FindLongURL(shortPath string) (string, error)
	FindURLsByUser(userID string) (map[string]string, error)
	InsertNewURLPair(userID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
}
