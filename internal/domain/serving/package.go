package serving

import "kite/internal/domain/shared"

type Package struct {
	ID      shared.PackageID
	Version Version
	Files   []PackageFile
}
