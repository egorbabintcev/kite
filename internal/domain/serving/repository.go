package serving

import "kite/internal/domain/shared"

type PackageRepository interface {
	Get(id shared.PackageID, version Version) (*Package, error)
}
