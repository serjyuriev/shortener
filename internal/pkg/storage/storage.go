package storage

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Store interface {
	FindLongURL(sp string) (string, error)
	InsertNewURLPair(sp string, l string) error
}

type store struct {
	URLs            map[string]string
	fileStoragePath string
	useFileStorage  bool
}

var (
	ErrNoURLWasFound = errors.New("no URL was found")
)

func NewStore(fileStoragePath string) (Store, error) {
	s := store{
		fileStoragePath: fileStoragePath,
		useFileStorage:  fileStoragePath != "",
	}
	if s.useFileStorage {
		if err := s.loadDataFromFile(); err != nil {
			return nil, err
		}
	} else {
		s.URLs = make(map[string]string)
	}
	return &s, nil
}

func (s *store) FindLongURL(sp string) (string, error) {
	l, ok := s.URLs[sp]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l, nil
}

func (s *store) InsertNewURLPair(sp string, l string) error {
	s.URLs[sp] = l
	if err := s.writeDataToFile(); err != nil {
		delete(s.URLs, sp)
		return err
	}
	return nil
}

func (s *store) writeDataToFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(s.URLs)
	if err != nil {
		return err
	}
	num, err := file.Write(data)
	if err != nil {
		return err
	}
	log.Printf("Number of bytes written: %d", num)
	return nil
}

func (s *store) loadDataFromFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	s.URLs = make(map[string]string)
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	if err = json.Unmarshal(b, &s.URLs); err != nil {
		return err
	}
	return nil
}
