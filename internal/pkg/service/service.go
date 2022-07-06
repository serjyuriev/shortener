package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

// Job contains information about worker's job.
type Job struct {
	Ctx    context.Context
	UserID string
	URLs   []string
}

// Stats contains amount of users and urls in the system.
type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// Service provides method of application service layer.
type Service interface {
	DeleteURLs(userID string, urls []string)
	FindByOriginalURL(ctx context.Context, originalURL string) (string, error)
	FindOriginalURL(ctx context.Context, shortPath string) (string, error)
	FindURLsByUser(ctx context.Context, userID string) (map[string]string, error)
	InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error
	InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
	GetStats(ctx context.Context) (*Stats, error)
}

type service struct {
	jobChan chan *Job
	store   storage.Store
}

// NewService initializes application service layer.
func NewService() (Service, error) {
	cfg := config.GetConfig()

	var s storage.Store
	var err error
	if cfg.DatabaseDSN != "" {
		s, err = storage.NewPgStore(cfg.DatabaseDSN)
	} else {
		s, err = storage.NewFileStore(cfg.FileStoragePath)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to create new storage:\n%w", err)
	}

	svc := &service{
		store:   s,
		jobChan: make(chan *Job),
	}

	for i := 0; i < 5; i++ {
		go func() {
			for job := range svc.jobChan {
				svc.deleteURLs(job.Ctx, job.UserID, job.URLs)
			}
		}()
	}

	return svc, nil
}

// DeleteURLs creates a job for removing URLs from storage and sends it into job channel.
func (s *service) DeleteURLs(userID string, urls []string) {
	s.jobChan <- &Job{
		Ctx:    context.Background(),
		UserID: userID,
		URLs:   urls,
	}
}

// FindByOriginalURL searches for short URL with corresponding original URL in application storage.
func (s *service) FindByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	shorty, err := s.store.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("unable to find short url:\n%w", err)
	}
	return shorty, nil
}

// FindOriginalURL searches for original URL with corresponding short URL in application storage.
func (s *service) FindOriginalURL(ctx context.Context, shortPath string) (string, error) {
	original, err := s.store.FindOriginalURL(ctx, shortPath)
	if err != nil {
		return "", fmt.Errorf("unable to find original url:\n%w", err)
	}
	return original, nil
}

// FindURLsByUser returns all URLs from application storage that were added by user with provided ID.
func (s *service) FindURLsByUser(ctx context.Context, userID string) (map[string]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user id:\n%w", err)
	}

	m, err := s.store.FindURLsByUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unable to find urls by user:\n%w", err)
	}
	return m, nil
}

// InsertManyURLs inserts provided short URL - original URL pairs into application storage.
func (s *service) InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse user id:\n%w", err)
	}

	if err = s.store.InsertManyURLs(ctx, uid, urls); err != nil {
		return fmt.Errorf("unable to insert many urls:\n%w", err)
	}
	return nil
}

// InsertNewURLPair inserts provided short URL - original URL pair into application storage.
func (s *service) InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse user id:\n%w", err)
	}

	if err = s.store.InsertNewURLPair(ctx, uid, shortPath, originalURL); err != nil {
		return fmt.Errorf("unable to insert url pair:\n%w", err)
	}
	return nil
}

// Ping performs a healthcheck of application storage.
func (s *service) Ping(ctx context.Context) error {
	if err := s.store.Ping(ctx); err != nil {
		return fmt.Errorf("unable to perform ping:\n%w", err)
	}
	return nil
}

// GetStats returns amount of users and URLs in the system.
func (s *service) GetStats(ctx context.Context) (*Stats, error) {
	urls, err := s.store.CountURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to count urls:\n%w", err)
	}
	users, err := s.store.CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to count users:\n%w", err)
	}

	return &Stats{URLs: urls, Users: users}, nil
}

func (s *service) deleteURLs(ctx context.Context, userID string, urls []string) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unable to parse user id (%s): %v", userID, err)
		return
	}

	if err = s.store.DeleteManyURLs(ctx, uid, urls); err != nil {
		log.Printf("unable to delete urls: %v", err)
	}
}
