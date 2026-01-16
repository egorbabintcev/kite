package application

import "context"

type GetPackageArtifact interface {
	Execute(context.Context, GetPackageArtifactRequest) (*GetPackageArtifactResponse, error)
}
