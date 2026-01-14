package resolution

import (
	"github.com/Masterminds/semver/v3"
)

type VersionRange struct {
	raw        string
	constraint *semver.Constraints
}

func NewVersionRange(raw string) (VersionRange, error) {
	if raw == "" {
		return VersionRange{}, ErrVersionRangeEmpty
	}

	constraint, err := semver.NewConstraint(raw)
	if err != nil {
		return VersionRange{}, ErrVersionRangeInvalid
	}

	return VersionRange{
		raw:        raw,
		constraint: constraint,
	}, nil
}

func (vr VersionRange) String() string {
	return vr.constraint.String()
}

func (vr VersionRange) Match(version Version) bool {
	return vr.constraint.Check(version.version)
}

func (vr VersionRange) Equal(other VersionRange) bool {
	return vr.String() == other.String()
}
