package application

import (
	"io"
	"time"
)

type GetPackageArtifactRequest struct {
	Name         string
	Scope        string
	VersionQuery string
	Path         string
}

type ArtifactServe struct {
	Stream  io.ReadSeekCloser
	Name    string
	ModTime time.Time
}

type ArtifactRedirect struct {
	Scope   string
	Name    string
	Version string
	Path    string
}

type GetPackageArtifactResponse struct {
	Serve    *ArtifactServe
	Redirect *ArtifactRedirect
}
