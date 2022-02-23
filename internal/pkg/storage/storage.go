package storage

import "errors"

type ShortPath string
type LongURL string

type Store interface {
	FindLongURL(sp ShortPath) (LongURL, error)
	InsertNewURLPair(sp ShortPath, l LongURL) error
}

var (
	ErrNoURLWasFound = errors.New("no URL was found")
)

func NewStore(fileStoragePath string) (Store, error) {
	if fileStoragePath == "" {
		return newInmemStore(), nil
	} else {
		return newFileStore(fileStoragePath)
	}
}
