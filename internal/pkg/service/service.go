package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
)

type Service interface {
	FindByOriginalURL(ctx context.Context, originalURL string) (string, error)
	FindOriginalURL(ctx context.Context, shortPath string) (string, error)
	FindURLsByUser(ctx context.Context, userID string) (map[string]string, error)
	InsertManyURLs(ctx context.Context, userID string, urls map[string]string) error
	InsertNewURLPair(ctx context.Context, userID, shortPath, originalURL string) error
	Ping(ctx context.Context) error
}

type service struct {
	store storage.Store
}

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
	return &service{store: s}, nil
}

func (s *service) FindByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	shorty, err := s.store.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("unable to find short url:\n%w", err)
	}
	return shorty, nil
}

func (s *service) FindOriginalURL(ctx context.Context, shortPath string) (string, error) {
	original, err := s.store.FindOriginalURL(ctx, shortPath)
	if err != nil {
		return "", fmt.Errorf("unable to find original url:\n%w", err)
	}
	return original, nil
}

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

func (s *service) Ping(ctx context.Context) error {
	if err := s.store.Ping(ctx); err != nil {
		return fmt.Errorf("unable to perform ping:\n%w", err)
	}
	return nil
}
