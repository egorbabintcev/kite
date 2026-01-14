package serving

import (
	"github.com/Masterminds/semver/v3"
)

type Version struct {
	value string
}

func NewVersion(raw string) (Version, error) {
	if raw == "" {
		return Version{}, ErrVersionEmpty
	}

	v, err := semver.StrictNewVersion(raw)
	if err != nil {
		return Version{}, ErrVersionInvalid
	}

	return Version{value: v.String()}, nil
}

func (v Version) String() string {
	return v.value
}
