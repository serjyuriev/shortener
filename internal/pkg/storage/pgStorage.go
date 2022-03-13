package storage

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type pgStore struct {
	db *sql.DB
	cs string
}

func NewPgStore(connectionString string) (Store, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("unable to open database: %v\n", err)
		return nil, err
	}
	s := &pgStore{
		cs: connectionString,
		db: db,
	}
	return s, nil
}

func (s *pgStore) FindLongURL(shortPath string) (string, error) {
	return "", ErrNotImplementedYet
}

func (s *pgStore) FindURLsByUser(userID string) (map[string]string, error) {
	return nil, ErrNotImplementedYet
}

func (s *pgStore) InsertNewURLPair(userID, shortPath, originalURL string) error {
	return ErrNotImplementedYet
}

func (s *pgStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
