package storage

type inMemStore struct {
	URLs map[ShortPath]LongURL
}

func newInmemStore() *inMemStore {
	s := inMemStore{
		URLs: make(map[ShortPath]LongURL),
	}
	return &s
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
