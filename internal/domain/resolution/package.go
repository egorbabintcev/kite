package resolution

import (
	"fmt"
	"kite/internal/domain/shared"
	"slices"
)

type Package struct {
	ID       shared.PackageID
	Versions []Version
	Tags     map[VersionTag]Version
}

func NewPackage(ID shared.PackageID, versions []Version, tags map[VersionTag]Version) (*Package, error) {
	if len(versions) == 0 {
		return nil, ErrPackageNoVersions
	}

	hasLatest := false
	for tag := range tags {
		if tag.IsLatest() {
			hasLatest = true
			break
		}
	}
	if !hasLatest {
		return nil, ErrPackageNoLatestTag
	}

	return &Package{
		ID:       ID,
		Versions: versions,
		Tags:     tags,
	}, nil
}

func (p *Package) ResolveVersionQuery(query VersionQuery) (Version, error) {
	if query.Tag != nil {
		if version, ok := p.Tags[*query.Tag]; ok {
			return version, nil
		} else {
			return Version{}, ErrVersionNotFound
		}
	}

	if query.Range != nil {
		rng := *query.Range

		v, _ := NewVersion("20.0.0")
		fmt.Printf("%v", rng.Match(v))
		candidates := make([]Version, 0)

		for _, version := range p.Versions {
			if rng.Match(version) {
				candidates = append(candidates, version)
			}
		}

		if len(candidates) == 0 {
			return Version{}, ErrVersionNotFound
		}

		slices.SortFunc(candidates, func(a Version, b Version) int {
			return a.Compare(b)
		})
		return candidates[len(candidates)-1], nil
	}

	return Version{}, ErrVersionNotFound
}
