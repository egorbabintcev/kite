package serving

import (
	"time"
)

type PackageFile struct {
	Path    string
	ModTime time.Time
}

func NewPackageFile(path string, modTime time.Time) PackageFile {
	return PackageFile{
		Path:    path,
		ModTime: modTime,
	}
}
