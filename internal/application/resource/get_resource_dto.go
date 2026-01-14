package resource

import (
	"io"
	"time"
)

type GetResourceUCRequest struct {
	Name         string
	Scope        string
	VersionQuery string
	Path         string
}

type ResourceServe struct {
	Stream  io.ReadSeekCloser
	Name    string
	ModTime time.Time
}

type ResourceRedirect struct {
	Scope   string
	Name    string
	Version string
	Path    string
}

type GetResourceUCResponse struct {
	Serve    *ResourceServe
	Redirect *ResourceRedirect
}
