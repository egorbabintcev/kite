package resolution

import (
	"fmt"
	"kite/internal/domain/shared"
)

type Package struct {
	ID       shared.PackageID
	Versions []Version
	Tags     map[VersionTag]Version
}

func NewPackage(ID shared.PackageID, versions []Version, tags map[VersionTag]Version) (*Package, error) {
	if len(versions) == 0 {
		return nil, fmt.Errorf("package must have at least one version")
	}

	hasLatest := false
	for tag := range tags {
		if tag.IsLatestTag() {
			hasLatest = true
			break
		}
	}
	if !hasLatest {
		return nil, fmt.Errorf("package must have a latest tag")
	}

	return &Package{
		ID:       ID,
		Versions: versions,
		Tags:     tags,
	}, nil
}
