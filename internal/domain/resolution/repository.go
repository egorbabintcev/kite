package resolution

import "kite/internal/domain/shared"

type PackageRepository interface {
	Get(id shared.PackageID) (*Package, error)
}
