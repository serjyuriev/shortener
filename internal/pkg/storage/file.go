package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type keyValPair struct {
	Short string
	Long  string
}

type fileStore struct {
	URLs            map[string]string
	fileStoragePath string
}

func newFileStore(fileStoragePath string) (*fileStore, error) {
	s := fileStore{
		fileStoragePath: fileStoragePath,
	}
	err := s.loadDataFromFile()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *fileStore) FindLongURL(sp string) (string, error) {
	l, ok := s.URLs[sp]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l, nil
}

func (s *fileStore) InsertNewURLPair(sp string, l string) error {
	e := make(chan error, 1)
	go func() {
		e <- s.writeDataToFile(sp, l)
		close(e)
	}()

	s.URLs[sp] = l
	return <-e
}

func (s *fileStore) writeDataToFile(sp string, l string) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	p := keyValPair{Short: sp, Long: l}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	num, err := file.Write(data)
	if err != nil {
		return err
	}
	log.Printf("Number of bytes written: %d", num)
	return nil
}

func (s *fileStore) loadDataFromFile() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	s.URLs = make(map[string]string)
	for scanner.Scan() {
		var p keyValPair
		err := json.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			return err
		}
		s.URLs[p.Short] = p.Long
	}
	return scanner.Err()
}
