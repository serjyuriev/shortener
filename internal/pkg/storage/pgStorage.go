package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type pgStore struct {
	db *sql.DB
	cs string
}

func NewPgStore(connectionString string) (Store, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database:\n%w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)
	s := &pgStore{
		cs: connectionString,
		db: db,
	}

	if _, err = s.db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
			short_id TEXT PRIMARY KEY,
			original_url TEXT NOT NULL,
			added_by_user TEXT NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`); err != nil {
		return nil, fmt.Errorf("unable to execute create statements:\n%w", err)
	}

	return s, nil
}

func (s *pgStore) FindLongURL(ctx context.Context, shortPath string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_id = $1", shortPath)
	var long string
	row.Scan(&long)
	if row.Err() != nil {
		return "", fmt.Errorf("unable to execute query:\n%w", row.Err())
	}

	return long, nil
}

func (s *pgStore) FindURLsByUser(ctx context.Context, userID string) (map[string]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user uuid:\n%w", err)
	}

	rows, err := s.db.QueryContext(ctx, "SELECT short_id, original_url FROM urls WHERE added_by_user = $1", uid.String())
	if err != nil {
		return nil, fmt.Errorf("unable to execute query:\n%w", err)
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var short, long string
		if err := rows.Scan(&short, &long); err != nil {
			return nil, fmt.Errorf("unable to scan values:\n%w", err)
		}
		urls[short] = long
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to execute query:\n%w", err)
	}

	if len(urls) == 0 {
		return nil, ErrNoURLWasFound
	}

	return urls, nil
}

func (s *pgStore) InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse user uuid:\n%w", err)
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return fmt.Errorf("unable to begin transaction:\n%w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls(short_id, original_url, added_by_user) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("unable to prepare sql statement:\n%w", err)
	}
	defer stmt.Close()

	for short, long := range urls {
		if _, err = stmt.Exec(short, long, uid.String()); err != nil {
			return fmt.Errorf("unable to execute sql statement:\n%w", err)
		}
	}

	return tx.Commit()
}

func (s *pgStore) InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse user uuid:\n%w", err)
	}

	if _, err = s.db.ExecContext(ctx, "INSERT INTO urls VALUES ($1, $2, $3)",
		shortPath, originalURL, uid.String()); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrNotUniqueOriginalURL
			}
		}
		return fmt.Errorf("unable to insert values:\n%w", err)
	}

	return nil
}

func (s *pgStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *pgStore) FindByLongURL(ctx context.Context, long string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT short_id FROM urls WHERE original_url = $1", long)
	row.Scan(&long)
	if row.Err() != nil {
		return "", fmt.Errorf("unable to execute query:\n%w", row.Err())
	}

	return long, nil
}
