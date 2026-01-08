package serving

import "kite/internal/domain/shared_kernel"

type Package struct {
	ID shared_kernel.PackageID
	Files []*PackageFile
}