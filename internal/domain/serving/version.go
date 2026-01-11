package serving

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type Version struct {
	value string
}

func NewVersion(raw string) (Version, error) {
	v, err := semver.StrictNewVersion(raw)
	if err != nil {
		return Version{}, fmt.Errorf("failed to create version")
	}

	return Version{value: v.String()}, nil
}

func (v Version) String() string {
	return v.value
}
