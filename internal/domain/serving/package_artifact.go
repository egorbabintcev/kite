package serving

import (
	"io"
	"kite/internal/domain/shared"
	"path/filepath"
	"time"
)

type PackageArtifactPath = string

type PackageArtifactContent interface {
	io.ReadSeekCloser
}

type PackageArtifact struct {
	id      shared.PackageID
	version Version
	path    PackageArtifactPath
	content PackageArtifactContent
	modTime time.Time
}

func NewPackageArtifact(id shared.PackageID, version Version, path PackageArtifactPath, content PackageArtifactContent, modTime time.Time) (*PackageArtifact, error) {
	return &PackageArtifact{
		id:      id,
		version: version,
		path:    path,
		content: content,
		modTime: modTime,
	}, nil
}

func (pa *PackageArtifact) Name() string {
	return filepath.Base(pa.path)
}

func (pa *PackageArtifact) ModTime() time.Time {
	return pa.modTime
}

func (pa *PackageArtifact) Content() PackageArtifactContent {
	return pa.content
}
