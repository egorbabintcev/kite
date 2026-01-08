package resolution

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type VersionRange struct {
	raw        string
	constraint *semver.Constraints
}

func NewVersionRange(raw string) (*VersionRange, error) {
	constraint, err := semver.NewConstraint(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to create version range")
	}

	return &VersionRange{
		raw:        raw,
		constraint: constraint,
	}, nil
}

func (vr *VersionRange) String() string {
	return vr.constraint.String()
}

func (vr *VersionRange) Match(version *Version) bool {
	return vr.constraint.Check(version.version)
}
