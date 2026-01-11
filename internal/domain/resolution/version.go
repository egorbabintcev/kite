package resolution

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type Version struct {
	raw     string
	version *semver.Version
}

func NewVersion(raw string) (Version, error) {
	version, err := semver.StrictNewVersion(raw)
	if err != nil {
		return Version{}, fmt.Errorf("failed to create version")
	}

	return Version{
		raw:     raw,
		version: version,
	}, nil
}

func (v Version) String() string {
	return v.version.String()
}

func (v Version) Compare(other Version) int {
	return v.version.Compare(other.version)
}

func (v Version) Equal(other Version) bool {
	return v.version.Equal(other.version)
}
