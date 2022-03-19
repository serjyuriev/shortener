package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type link struct {
	Original string
	User     uuid.UUID
}

type fileStore struct {
	URLs            map[string]link
	fileStoragePath string
	useFileStorage  bool
}

func NewFileStore(fileStoragePath string) (Store, error) {
	s := &fileStore{
		fileStoragePath: fileStoragePath,
		useFileStorage:  fileStoragePath != "",
	}
	if s.useFileStorage {
		if err := s.loadDataFromFile(); err != nil {
			// log.Printf("unable to load data from file: %v\n", err)
			return nil, fmt.Errorf("unable to load data from file: %w", err)
		}
	} else {
		s.URLs = make(map[string]link)
	}
	return s, nil
}

func (s *fileStore) FindByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	for _, v := range s.URLs {
		if v.Original == originalURL {
			return v.Original, nil
		}
	}
	return "", ErrNoURLWasFound
}

func (s *fileStore) FindOriginalURL(ctx context.Context, shortPath string) (string, error) {
	l, ok := s.URLs[shortPath]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l.Original, nil
}

func (s *fileStore) FindURLsByUser(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	userURLs := make(map[string]string)
	for k, v := range s.URLs {
		if v.User == userID {
			userURLs[k] = v.Original
		}
	}
	if len(userURLs) == 0 {
		return nil, ErrNoURLWasFound
	}
	return userURLs, nil
}

func (s *fileStore) InsertManyURLs(ctx context.Context, userID uuid.UUID, urls map[string]string) error {
	oldMap := make(map[string]link)
	for v, k := range s.URLs {
		oldMap[v] = k
	}

	for short, long := range urls {
		newLink := link{
			Original: long,
			User:     userID,
		}
		s.URLs[short] = newLink
	}
	if s.useFileStorage {
		if err := s.writeDataToFile(); err != nil {
			s.URLs = oldMap
			return err
		}
	}
	return nil
}

func (s *fileStore) InsertNewURLPair(ctx context.Context, userID uuid.UUID, shortPath, originalURL string) error {
	newLink := link{
		Original: originalURL,
		User:     userID,
	}
	s.URLs[shortPath] = newLink
	if s.useFileStorage {
		if err := s.writeDataToFile(); err != nil {
			delete(s.URLs, shortPath)
			return err
		}
	}
	return nil
}

func (s *fileStore) Ping(ctx context.Context) error {
	return nil
}

func (s *fileStore) loadDataFromFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("unable to open file %s: %v\n", s.fileStoragePath, err)
		return err
	}
	defer file.Close()

	s.URLs = make(map[string]link)
	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("unable to read from file: %v\n", err)
		return err
	}
	if len(b) == 0 {
		return nil
	}
	if err = json.Unmarshal(b, &s.URLs); err != nil {
		log.Printf("unable to unmarshal json: %v\n", err)
		return err
	}
	return nil
}

func (s *fileStore) writeDataToFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("unable to open file %s: %v\n", s.fileStoragePath, err)
		return err
	}
	defer file.Close()

	data, err := json.Marshal(s.URLs)
	if err != nil {
		log.Printf("unable to marshal map to json: %v\n", err)
		return err
	}
	num, err := file.Write(data)
	if err != nil {
		log.Printf("unable to write data to file: %v\n", err)
		return err
	}
	log.Printf("Number of bytes written: %d", num)
	return nil
}
