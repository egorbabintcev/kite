package resource

import (
	"context"
	"kite/internal/domain/resolution"
	"kite/internal/domain/serving"
	"kite/internal/domain/shared"
)

type GetResourceUC struct {
	resolutionRepo resolution.PackageRepository
	servingRepo    serving.PackageArtifactRepository
}

func NewGetResourceUC(resolutionRepo resolution.PackageRepository, servingRepo serving.PackageArtifactRepository) *GetResourceUC {
	return &GetResourceUC{
		resolutionRepo: resolutionRepo,
		servingRepo:    servingRepo,
	}
}

func (uc *GetResourceUC) Execute(ctx context.Context, r GetResourceUCRequest) (*GetResourceUCResponse, error) {
	id, err := shared.NewPackageID(r.Scope, r.Name)
	if err != nil {
		return nil, err
	}

	pkg, err := uc.resolutionRepo.Get(ctx, id)
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
		return &GetResourceUCResponse{
			Redirect: &ResourceRedirect{
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

	artifact, err := uc.servingRepo.Get(ctx, id, v, r.Path)
	if err != nil {
		return nil, err
	}

	return &GetResourceUCResponse{
		Serve: &ResourceServe{
			Name:    artifact.Name(),
			ModTime: artifact.ModTime(),
			Stream:  artifact.Content(),
		},
	}, nil
}
