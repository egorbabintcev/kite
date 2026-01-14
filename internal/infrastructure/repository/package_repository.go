package repository

import (
	"context"
	"kite/internal/domain/resolution"
	"kite/internal/domain/shared"
	"kite/internal/infrastructure/registry"
)

type PackageRepository struct {
	registry *registry.HttpClient
}

func NewPackageRepository(r *registry.HttpClient) *PackageRepository {
	return &PackageRepository{registry: r}
}

func (pr *PackageRepository) Get(ctx context.Context, id shared.PackageID) (*resolution.Package, error) {
	meta, err := pr.registry.FetchMetadata(ctx, id.String())
	if err != nil {
		return nil, resolution.ErrPackageNotFound
	}

	versions := make([]resolution.Version, 0, len(meta.Metadata.Versions))
	for _, rawVersion := range meta.Metadata.Versions {
		v, err := resolution.NewVersion(rawVersion)
		if err != nil {
			return nil, err
		}

		versions = append(versions, v)
	}

	tags := make(map[resolution.VersionTag]resolution.Version, len(meta.Metadata.Tags))
	for tag, version := range meta.Metadata.Tags {
		v, err := resolution.NewVersion(version)
		if err != nil {
			return nil, err
		}

		t, err := resolution.NewVersionTag(tag)
		if err != nil {
			return nil, err
		}

		tags[t] = v
	}

	pkg, err := resolution.NewPackage(id, versions, tags)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
