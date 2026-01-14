package serving

import (
	"context"
	"kite/internal/domain/shared"
)

type PackageArtifactRepository interface {
	Get(ctx context.Context, id shared.PackageID, version Version, path PackageArtifactPath) (*PackageArtifact, error)
}
