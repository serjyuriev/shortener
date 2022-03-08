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
	FindLongURL(userID, shortPath string) (string, error)
	FindURLsByUser(userID string) (map[string]string, error)
	InsertNewURLPair(userID, shortPath, l string) error
	IsUserExists(uid uuid.UUID) bool
}

type store struct {
	URLs            map[uuid.UUID]map[string]string
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
		s.URLs = make(map[uuid.UUID]map[string]string)
	}
	return &s, nil
}

func (s *store) FindLongURL(userID, shortPath string) (string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}
	l, ok := s.URLs[uid][shortPath]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l, nil
}

func (s *store) FindURLsByUser(userID string) (map[string]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	if userURLs, ok := s.URLs[uid]; ok {
		return userURLs, nil
	}
	return nil, ErrNoURLWasFound
}

func (s *store) InsertNewURLPair(userID, shortPath, l string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if !s.IsUserExists(uid) {
		s.URLs[uid] = make(map[string]string)
	}
	s.URLs[uid][shortPath] = l
	if s.useFileStorage {
		if err := s.writeDataToFile(); err != nil {
			delete(s.URLs[uid], shortPath)
			return err
		}
	}
	return nil
}

func (s *store) IsUserExists(uid uuid.UUID) bool {
	_, ok := s.URLs[uid]
	return ok
}

func (s *store) loadDataFromFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("unable to open file %s: %v\n", s.fileStoragePath, err)
		return err
	}
	defer file.Close()

	s.URLs = make(map[uuid.UUID]map[string]string)
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
