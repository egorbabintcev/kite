package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
)

var (
	ErrInternal = errors.New("internal error")
)

type Service struct {
	logger   *slog.Logger
	cache    Cache
	registry RegistryClient
}

func NewService(l *slog.Logger, c Cache, r RegistryClient) *Service {
	l = l.With(slog.String("component", "core"))

	return &Service{
		logger:   l,
		cache:    c,
		registry: r,
	}
}

func (s *Service) GetResource(ctx context.Context, scope, name, version, path string) (*GetResourceResponse, error) {
	fullName := name
	if scope != "" {
		fullName = fmt.Sprintf("@%s/%s", scope, name)
	}

	cacheKey := filepath.Join(fullName, version)

	if exists := s.cache.Exists(cacheKey); !exists {
		res, err := s.registry.FetchPackage(ctx, scope, name, version)
		if err != nil {
			return nil, fmt.Errorf("error fetching resource: %w", err)
		}

		for _, f := range res.Files {
			if err := s.cache.Put(filepath.Join(cacheKey, f.Path), f.Content); err != nil {
				return nil, fmt.Errorf("error putting resource to cache: %w", err)
			}
		}
	}

	rsc, info, err := s.cache.Get(filepath.Join(cacheKey, path))
	if err != nil {
		return nil, fmt.Errorf("error getting resource from cache: %w", err)
	}

	return &GetResourceResponse{
		Resource: Resource{
			Stream:  rsc,
			Name:    info.Name(),
			ModTime: info.ModTime(),
		},
	}, nil
}
