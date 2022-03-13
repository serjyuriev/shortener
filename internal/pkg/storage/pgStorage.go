package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
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
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)
	s := &pgStore{
		cs: connectionString,
		db: db,
	}

	if _, err = s.db.Exec("CREATE TABLE IF NOT EXISTS urls ( short_id TEXT PRIMARY KEY, original_url TEXT NOT NULL, added_by_user TEXT NOT NULL );"); err != nil {
		log.Printf("unable to execute create statement: %v\n", err)
		return nil, err
	}

	return s, nil
}

func (s *pgStore) FindLongURL(ctx context.Context, shortPath string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_id = $1", shortPath)
	var long string
	row.Scan(&long)
	if row.Err() != nil {
		log.Printf("unable to scan value: %v\n", row.Err())
		return "", row.Err()
	}

	return long, nil
}

func (s *pgStore) FindURLsByUser(ctx context.Context, userID string) (map[string]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, "SELECT short_id, original_url FROM urls WHERE added_by_user = $1", uid.String())
	if err != nil {
		log.Printf("unable to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var short, long string
		if err := rows.Scan(&short, &long); err != nil {
			log.Printf("unable to scan values: %v\n", err)
			return nil, err
		}
		urls[short] = long
	}

	if rows.Err() != nil {
		log.Printf("unable to scan values: %v\n", err)
		return nil, rows.Err()
	}

	if len(urls) == 0 {
		return nil, ErrNoURLWasFound
	}

	return urls, nil
}

func (s *pgStore) InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unable to parse user uuid: %v\n", err)
		return err
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		log.Printf("unable to start transaction: %v\n", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls(short_id, original_url, added_by_user) VALUES ($1, $2, $3)")
	if err != nil {
		log.Printf("unable to prepare sql statement: %v\n", err)
		return err
	}
	defer stmt.Close()

	for short, long := range urls {
		if _, err = stmt.Exec(short, long, uid.String()); err != nil {
			log.Printf("unable to execute sql statement: %v\n", err)
			return err
		}
	}

	return tx.Commit()
}

func (s *pgStore) InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	if _, err = s.db.ExecContext(ctx, "INSERT INTO urls VALUES ($1, $2, $3)",
		shortPath, originalURL, uid.String()); err != nil {
		log.Printf("unable to insert values: %v\n", err)
		return err
	}

	return nil
}

func (s *pgStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
