package core

import (
	"context"
	"io"
	"kite/internal/infrastructure/registry"
	"os"
)

type Cache interface {
	Get(path string) (io.ReadSeekCloser, os.FileInfo, error)
	Put(path string, content []byte) error
	Exists(path string) bool
}

type RegistryClient interface {
	FetchPackage(ctx context.Context, scope, name, version string) (*registry.GetPackageResponse, error)
	FetchMetadata(ctx context.Context, scope, name string) (*registry.GetMetadataResponse, error)
}
