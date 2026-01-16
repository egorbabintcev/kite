package application

import (
	"context"
	"kite/internal/domain/resolution"
	"kite/internal/domain/serving"
	"kite/internal/domain/shared"
)

type GetPackageArtifactUC struct {
	packageRepo  resolution.PackageRepository
	artifactRepo serving.PackageArtifactRepository
}

func NewGetPackageArtifactUC(pr resolution.PackageRepository, ar serving.PackageArtifactRepository) *GetPackageArtifactUC {
	return &GetPackageArtifactUC{
		packageRepo:  pr,
		artifactRepo: ar,
	}
}

func (uc *GetPackageArtifactUC) Execute(ctx context.Context, r GetPackageArtifactRequest) (*GetPackageArtifactResponse, error) {
	id, err := shared.NewPackageID(r.Scope, r.Name)
	if err != nil {
		return nil, err
	}

	pkg, err := uc.packageRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	query, err := resolution.NewVersionQuery(r.VersionQuery)
	if err != nil {
		return nil, err
	}

	version, err := pkg.ResolveVersionQuery(query)
	if err != nil {
		return nil, err
	}

	if version.String() != r.VersionQuery {
		return &GetPackageArtifactResponse{
			Redirect: &ArtifactRedirect{
				Name:    id.Name(),
				Scope:   id.Scope(),
				Version: version.String(),
				Path:    r.Path,
			},
		}, nil
	}

	v, err := serving.NewVersion(version.String())
	if err != nil {
		return nil, err
	}

	artifact, err := uc.artifactRepo.Get(ctx, id, v, r.Path)
	if err != nil {
		return nil, err
	}

	return &GetPackageArtifactResponse{
		Serve: &ArtifactServe{
			Name:    artifact.Name(),
			ModTime: artifact.ModTime(),
			Stream:  artifact.Content(),
		},
	}, nil
}
