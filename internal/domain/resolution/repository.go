package resolution

import (
	"context"
	"kite/internal/domain/shared"
)

type PackageRepository interface {
	Get(ctx context.Context, id shared.PackageID) (*Package, error)
}
