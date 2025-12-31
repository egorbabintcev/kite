package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
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

	meta, err := s.registry.FetchMetadata(ctx, scope, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching resource metadata: %w", err)
	}

	if version == "" {
		version = "latest"
	}

	resolvedVersion, err := s.resolveVersion(version, meta.Metadata.Versions, meta.Metadata.Tags)
	if err != nil {
		return nil, fmt.Errorf("error resolving version: %w", err)
	}

	if resolvedVersion != version {
		return &GetResourceResponse{
			Redirect: &ResourceRedirect{
				Scope:   scope,
				Name:    name,
				Version: resolvedVersion,
				Path:    path,
			},
		}, nil
	}

	cacheKey := filepath.Join(fullName, resolvedVersion)

	if exists := s.cache.Exists(cacheKey); !exists {
		res, err := s.registry.FetchPackage(ctx, scope, name, resolvedVersion)
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
		Serve: &ResourceServe{
			Stream:  rsc,
			Name:    info.Name(),
			ModTime: info.ModTime(),
		},
	}, nil
}

func (s *Service) resolveVersion(version string, versions []string, tags map[string]string) (string, error) {
	if v, ok := tags[version]; ok {
		return v, nil
	}

	versionSet := make(map[string]struct{}, len(versions))
	for _, v := range versions {
		versionSet[v] = struct{}{}
	}

	if _, ok := versionSet[version]; ok {
		return version, nil
	}

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return "", fmt.Errorf("failed to create constraint: %w", err)
	}

	matched := make([]*semver.Version, 0)
	for _, v := range versions {
		sv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}

		if constraint.Check(sv) {
			matched = append(matched, sv)
		}
	}

	if len(matched) == 0 {
		return "", fmt.Errorf("failed to match constraint")
	}

	sort.Sort(semver.Collection(matched))
	return matched[len(matched)-1].Original(), nil
}
