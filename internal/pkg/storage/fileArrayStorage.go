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

type arrayLink struct {
	link
	Shortened string
}

type fileArrayStore struct {
	URLs            []arrayLink
	fileStoragePath string
	useFileStorage  bool
}

// NewFileStore initializes file storage.
func NewFileArrayStore(fileStoragePath string) (Store, error) {
	s := &fileArrayStore{
		fileStoragePath: fileStoragePath,
		useFileStorage:  fileStoragePath != "",
	}
	if s.useFileStorage {
		if err := s.loadDataFromFile(); err != nil {
			return nil, fmt.Errorf("unable to load data from file: %w", err)
		}
	} else {
		s.URLs = make([]arrayLink, 0)
	}
	return s, nil
}

// DeleteManyURLs removes provided URLs from file.
func (s *fileArrayStore) DeleteManyURLs(ctx context.Context, userID uuid.UUID, urls []string) error {
	return ErrNotImplementedYet
}

// FindByOriginalURL searches for short URL with corresponding original URL.
func (s *fileArrayStore) FindByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	for _, v := range s.URLs {
		if v.Original == originalURL {
			return v.Shortened, nil
		}
	}
	return "", ErrNoURLWasFound
}

// FindOriginalURL searches for original URL with corresponding short URL.
func (s *fileArrayStore) FindOriginalURL(ctx context.Context, shortPath string) (string, error) {
	for _, v := range s.URLs {
		if v.Shortened == shortPath {
			return v.Original, nil
		}
	}
	return "", ErrNoURLWasFound
}

// FindURLsByUser returns all URLs from application storage that were added by user with provided ID.
func (s *fileArrayStore) FindURLsByUser(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	userURLs := make(map[string]string)
	for _, v := range s.URLs {
		if v.User == userID {
			userURLs[v.Shortened] = v.Original
		}
	}
	if len(userURLs) == 0 {
		return nil, ErrNoURLWasFound
	}
	return userURLs, nil
}

// InsertManyURLs writes provided short URL - original URL pairs into a file.
func (s *fileArrayStore) InsertManyURLs(ctx context.Context, userID uuid.UUID, urls map[string]string) error {
	links := make([]arrayLink, 0)
	for k, v := range urls {
		links = append(
			links,
			arrayLink{
				Shortened: k,
				link: link{
					Original: v,
					User:     userID,
				},
			},
		)
	}
	oldURLs := s.URLs
	s.URLs = append(s.URLs, links...)
	if s.useFileStorage {
		if err := s.writeDataToFile(); err != nil {
			s.URLs = oldURLs
			return err
		}
	}
	return nil
}

// InsertNewURLPair writes provided short URL - original URL pair into a file.
func (s *fileArrayStore) InsertNewURLPair(ctx context.Context, userID uuid.UUID, shortPath, originalURL string) error {
	newLink := arrayLink{
		Shortened: shortPath,
		link: link{
			Original: originalURL,
			User:     userID,
		},
	}
	s.URLs = append(s.URLs, newLink)
	if s.useFileStorage {
		if err := s.writeDataToFile(); err != nil {
			if len(s.URLs) == 1 {
				s.URLs = make([]arrayLink, 0)
			} else {
				s.URLs = s.URLs[:len(s.URLs)-2]
			}
			return err
		}
	}
	return nil
}

// Ping does nothing.
func (s *fileArrayStore) Ping(ctx context.Context) error {
	return nil
}

func (s *fileArrayStore) loadDataFromFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("unable to open file %s: %v\n", s.fileStoragePath, err)
		return err
	}
	defer file.Close()

	s.URLs = make([]arrayLink, 0)
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

func (s *fileArrayStore) writeDataToFile() error {
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
