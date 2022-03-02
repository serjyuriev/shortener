package storage

type inMemStore struct {
	URLs map[string]string
}

func newInmemStore() *inMemStore {
	s := inMemStore{
		URLs: make(map[string]string),
	}
	return &s
}

func (s *inMemStore) FindLongURL(sp string) (string, error) {
	l, ok := s.URLs[sp]
	if !ok {
		return "", ErrNoURLWasFound
	}
	return l, nil
}

func (s *inMemStore) InsertNewURLPair(sp string, l string) error {
	s.URLs[sp] = l
	return nil
}
