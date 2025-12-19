package core

import (
	"io"
	"time"
)

type Resource struct {
	Stream  io.ReadSeekCloser
	Name    string
	ModTime time.Time
}

type GetResourceResponse struct {
	Resource Resource
}
