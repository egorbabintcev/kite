package repository

import (
	"context"
	"fmt"
	"kite/internal/domain/serving"
	"kite/internal/domain/shared"
	"kite/internal/infrastructure/cache"
	"kite/internal/infrastructure/registry"
	"path/filepath"
)

type PackageArtifactRepository struct {
	registry *registry.HttpClient
	cache    *cache.FS
}

func NewPackageArtifactRepository(registry *registry.HttpClient, cache *cache.FS) *PackageArtifactRepository {
	return &PackageArtifactRepository{
		registry: registry,
		cache:    cache,
	}
}

func (r *PackageArtifactRepository) Get(ctx context.Context, id shared.PackageID, version serving.Version, path serving.PackageArtifactPath) (*serving.PackageArtifact, error) {
	cacheKey := filepath.Join(id.String(), version.String())

	if !r.cache.Exists(cacheKey) {
		res, err := r.registry.FetchPackage(ctx, id.Scope(), id.Name(), version.String())
		if err != nil {
			return nil, serving.ErrArtifactNotFound
		}

		for _, f := range res.Files {
			if err := r.cache.Put(filepath.Join(cacheKey, f.Path), f.Content); err != nil {
				return nil, fmt.Errorf("error putting resource to cache: %w", err)
			}
		}
	}

	rsc, info, err := r.cache.Get(filepath.Join(cacheKey, path))
	if err != nil {
		return nil, fmt.Errorf("error getting resource from cache: %w", err)
	}

	artifact, err := serving.NewPackageArtifact(id, version, path, rsc, info.ModTime())
	if err != nil {
		return nil, err
	}

	return artifact, nil
}
