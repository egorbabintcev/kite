package resolution

import "kite/internal/domain/shared_kernel"

type PackageRepository interface {
	Get(id shared_kernel.PackageID) (*Package, error)
}
