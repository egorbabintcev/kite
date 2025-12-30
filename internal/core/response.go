package core

import (
	"io"
	"time"
)

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

type GetResourceResponse struct {
	Serve    *ResourceServe
	Redirect *ResourceRedirect
}
