package grpcsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/service"
	"github.com/serjyuriev/shortener/internal/pkg/shorty"
	"github.com/serjyuriev/shortener/internal/pkg/storage"
	g "github.com/serjyuriev/shortener/proto/grpchandlers"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handlers store link to apps config and service layer.
type Service struct {
	g.UnimplementedShortenerServer

	cfg       *config.Config
	shortySvc service.Service
}

// MakeService initializes application grpc service layer.
func MakeService() (*Service, error) {
	svc, err := service.NewService()
	if err != nil {
		return nil, fmt.Errorf("unable to create new service:\n%w", err)
	}

	return &Service{
		cfg:       config.GetConfig(),
		shortySvc: svc,
	}, nil
}

// DeleteURLs removes URLs provided by user from storage.
func (s *Service) DeleteURLs(ctx context.Context, in *g.DeleteURLsRequest) (*g.DeleteURLsResponse, error) {
	var res g.DeleteURLsResponse

	s.shortySvc.DeleteURLs(in.UserID, in.Originals)

	return &res, nil
}

// GetURL searches service store for provided short URL
// and, if such URL is found, sends a response with original URL.
func (s *Service) GetURL(ctx context.Context, in *g.GetURLRequest) (*g.GetURLResponse, error) {
	var res g.GetURLResponse

	original, err := s.shortySvc.FindOriginalURL(ctx, in.Shorty)
	if err != nil {
		if errors.Is(err, storage.ErrShortenedDeleted) {
			res.Error = fmt.Sprintf("link with path %s is deleted", in.Shorty)
		} else {
			res.Error = fmt.Sprintf("unable to find full URL for %s", in.Shorty)
		}
	} else {
		res.Original = original
	}
	return &res, nil
}

// GetUserURLs returns all URLs that were added by current user.
func (s *Service) GetUserURLs(ctx context.Context, in *g.GetUserURLsRequest) (*g.GetUserURLsResponse, error) {
	var res g.GetUserURLsResponse

	m, err := s.shortySvc.FindURLsByUser(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrNoURLWasFound) {
			res.Error = fmt.Sprintf("user %s does not have any URLs", in.UserID)
		} else {
			res.Error = "unable to find URLs"
		}
	} else {
		res.Shorties = make([]string, len(m))
		res.Originals = make([]string, len(m))
		for key, val := range m {
			res.Shorties = append(
				res.Shorties,
				fmt.Sprintf("%s/%s", s.cfg.BaseURL, key),
			)
			res.Originals = append(res.Originals, val)
		}
	}
	return &res, nil
}

// Ping provides health status of application.
func (s *Service) Ping(ctx context.Context, in *emptypb.Empty) (*g.PingResponse, error) {
	var res g.PingResponse
	res.Error = s.shortySvc.Ping(ctx).Error()
	return &res, nil
}

// PostBatch adds URLs provided by user into storage,
// returning shortened URLs with corresponding correlation ID.
func (s *Service) PostBatch(ctx context.Context, in *g.PostBatchRequest) (*g.PostBatchResponse, error) {
	var res g.PostBatchResponse
	res.CorrelationID = make([]string, len(in.Originals))
	res.Shorties = make([]string, len(in.Originals))
	m := make(map[string]string)

	for i := 0; i < len(in.Originals); i++ {
		res.CorrelationID = append(res.CorrelationID, in.CorrelationID[i])
		res.Shorties = append(res.Shorties, shorty.GenerateShortPath())
		m[res.Shorties[i]] = in.Originals[i]
	}

	if err := s.shortySvc.InsertManyURLs(ctx, in.UserID, m); err != nil {
		res.Error = "unable to insert urls"
	}

	return &res, nil
}

// PostURL reads a long URL provided in request and
// creates a corresponding short URL, storing both in apps store.
func (s *Service) PostURL(ctx context.Context, in *g.PostURLRequest) (*g.PostURLResponse, error) {
	var res g.PostURLResponse

	sh := shorty.GenerateShortPath()

	if err := s.shortySvc.InsertNewURLPair(ctx, in.UserID, sh, in.Url); err != nil {
		res.Error = fmt.Sprintf("unable to save URL %s", in.Url)
	} else {
		res.Shorty = fmt.Sprintf("%s/%s", s.cfg.BaseURL, sh)
	}

	return &res, nil
}
