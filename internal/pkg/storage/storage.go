package storage

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

type Store interface {
	FindLongURL(shortPath string) (string, error)
	FindURLsByUser(userID string) (map[string]string, error)
	InsertNewURLPair(userID, shortPath, originalURL string) error
}

type link struct {
	Original string
	User     uuid.UUID
}

type store struct {
	URLs            map[string]link
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
		s.URLs = make(map[string]link)
	}
	return &s, nil
}

func (s *store) FindLongURL(shortPath string) (string, error) {
	l, ok := s.URLs[shortPath]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l.Original, nil
}

func (s *store) FindURLsByUser(userID string) (map[string]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	userURLs := make(map[string]string)
	for k, v := range s.URLs {
		if v.User == uid {
			userURLs[k] = v.Original
		}
	}
	if len(userURLs) == 0 {
		return nil, ErrNoURLWasFound
	}
	return userURLs, nil
}

func (s *store) InsertNewURLPair(userID, shortPath, originalURL string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	newLink := link{
		Original: originalURL,
		User:     uid,
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

func (s *store) loadDataFromFile() error {
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

func (s *store) writeDataToFile() error {
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
