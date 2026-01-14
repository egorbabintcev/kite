package serving

import "errors"

var (
	ErrVersionEmpty     = errors.New("version cannot be empty")
	ErrVersionInvalid   = errors.New("version is invalid")
	ErrArtifactNotFound = errors.New("artifact not found")
)
