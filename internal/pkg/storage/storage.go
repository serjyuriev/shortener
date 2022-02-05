package storage

import (
	"context"
	"errors"
	"time"
)

type ShortPath string
type LongURL string

type Store interface {
	FindLongURL(ctx context.Context, sp ShortPath) (LongURL, error)
	InsertNewURLPair(ctx context.Context, sp ShortPath, l LongURL) error
}

type inMemStore struct {
	URLs map[ShortPath]LongURL
}

var (
	ErrNoURLWasFound = errors.New("no URL was found")
)

func NewStore(ctx context.Context) (*inMemStore, error) {
	ctx2, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	s := inMemStore{
		URLs: make(map[ShortPath]LongURL),
	}
	<-ctx2.Done()
	return &s, nil
}

func (s *inMemStore) FindLongURL(ctx context.Context, sp ShortPath) (LongURL, error) {
	ctx2, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	l, ok := s.URLs[sp]
	if !ok {
		return "", ErrNoURLWasFound
	}
	<-ctx2.Done()
	return l, nil
}

func (s *inMemStore) InsertNewURLPair(ctx context.Context, sp ShortPath, l LongURL) error {
	ctx2, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	s.URLs[sp] = l
	<-ctx2.Done()
	return nil
}
