package resolution

import "errors"

var (
	ErrVersionTagEmpty     = errors.New("version tag cannot be empty")
	ErrVersionRangeEmpty   = errors.New("version range cannot be empty")
	ErrVersionRangeInvalid = errors.New("invalid version range")
	ErrVersionEmpty        = errors.New("version cannot be empty")
	ErrVersionInvalid      = errors.New("invalid version")
	ErrPackageNoVersions   = errors.New("package must have at least one version")
	ErrPackageNoLatestTag  = errors.New("package must have a latest tag")
	ErrVersionQueryEmpty   = errors.New("version query cannot be empty")
	ErrVersionQueryInvalid = errors.New("invalid version query")
	ErrVersionNotFound     = errors.New("version not found")
	ErrPackageNotFound     = errors.New("package not found")
)
