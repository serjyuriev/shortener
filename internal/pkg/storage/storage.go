package storage

import (
	"errors"
)

type ShortPath string
type LongURL string

type Store interface {
	FindLongURL(sp ShortPath) (LongURL, error)
	InsertNewURLPair(sp ShortPath, l LongURL) error
}

type inMemStore struct {
	URLs map[ShortPath]LongURL
}

var (
	ErrNoURLWasFound = errors.New("no URL was found")
)

func NewStore() (*inMemStore, error) {
	s := inMemStore{
		URLs: make(map[ShortPath]LongURL),
	}
	return &s, nil
}

func (s *inMemStore) FindLongURL(sp ShortPath) (LongURL, error) {
	l, ok := s.URLs[sp]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l, nil
}

func (s *inMemStore) InsertNewURLPair(sp ShortPath, l LongURL) error {
	s.URLs[sp] = l
	return nil
}
