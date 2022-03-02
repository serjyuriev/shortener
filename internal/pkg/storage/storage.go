package storage

import "errors"

type Store interface {
	FindLongURL(sp string) (string, error)
	InsertNewURLPair(sp string, l string) error
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
