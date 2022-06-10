package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

// NewPgStore initializes PostgreSQL storage.
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
			added_by_user TEXT NOT NULL,
			is_deleted BOOLEAN NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);`); err != nil {
		return nil, fmt.Errorf("unable to execute create statements:\n%w", err)
	}

	return s, nil
}

// DeleteManyURLs removes provided URLs from database.
func (s *pgStore) DeleteManyURLs(ctx context.Context, userID uuid.UUID, urls []string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return fmt.Errorf("unable to begin transaction:\n%w", err)
	}
	defer tx.Rollback()

	cmd := fmt.Sprintf(
		"UPDATE urls SET is_deleted = TRUE WHERE added_by_user = $1 AND short_id in ('%s');",
		strings.Join(urls, "','"),
	)

	stmt, err := tx.Prepare(cmd)
	if err != nil {
		return fmt.Errorf("unable to prepare sql statement:\n%w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(userID.String()); err != nil {
		return fmt.Errorf("unable to execute sql statement:\n%w", err)
	}

	return tx.Commit()
}

// FindByOriginalURL searches for short URL with corresponding original URL in database.
func (s *pgStore) FindByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT short_id, is_deleted FROM urls WHERE original_url = $1", originalURL)
	var isDeleted bool
	if err := row.Scan(&originalURL, &isDeleted); err != nil {
		return "", fmt.Errorf("unable to execute query:\n%w", row.Err())
	}
	if row.Err() != nil {
		return "", fmt.Errorf("unable to execute query:\n%w", row.Err())
	}

	if isDeleted {
		return "", ErrShortenedDeleted
	}

	return originalURL, nil
}

// FindOriginalURL searches for original URL with corresponding short URL in database.
func (s *pgStore) FindOriginalURL(ctx context.Context, shortPath string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM urls WHERE short_id = $1", shortPath)
	var long string
	var isDeleted bool
	row.Scan(&long, &isDeleted)
	if row.Err() != nil {
		return "", fmt.Errorf("unable to execute query:\n%w", row.Err())
	}

	if isDeleted {
		return "", ErrShortenedDeleted
	}

	return long, nil
}

// FindURLsByUser returns all URLs from application storage that were added by user with provided ID.
func (s *pgStore) FindURLsByUser(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT short_id, original_url FROM urls WHERE added_by_user = $1 AND is_deleted != TRUE", userID.String())
	if err != nil {
		return nil, fmt.Errorf("unable to execute query:\n%w", err)
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var short, long string
		if err = rows.Scan(&short, &long); err != nil {
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

// InsertManyURLs writes provided short URL - original URL pairs into database.
func (s *pgStore) InsertManyURLs(ctx context.Context, userID uuid.UUID, urls map[string]string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return fmt.Errorf("unable to begin transaction:\n%w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls(short_id, original_url, added_by_user, is_deleted) VALUES ($1, $2, $3, FALSE)")
	if err != nil {
		return fmt.Errorf("unable to prepare sql statement:\n%w", err)
	}
	defer stmt.Close()

	for short, long := range urls {
		if _, err = stmt.Exec(short, long, userID.String()); err != nil {
			return fmt.Errorf("unable to execute sql statement:\n%w", err)
		}
	}

	return tx.Commit()
}

// InsertNewURLPair writes provided short URL - original URL pair into database.
func (s *pgStore) InsertNewURLPair(ctx context.Context, userID uuid.UUID, shortPath, originalURL string) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO urls VALUES ($1, $2, $3, $4)",
		shortPath,
		originalURL,
		userID.String(),
		false,
	); err != nil {
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

// Ping checks connection with database.
func (s *pgStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
