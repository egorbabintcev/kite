package registry

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpClient struct {
	url string
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{url: url}
}

func (c *HttpClient) FetchPackage(ctx context.Context, scope, name, version string) (*GetPackageResponse, error) {
	fullName := name
	if scope != "" {
		fullName = fmt.Sprintf("@%s/%s", scope, name)
	}

	fetchResponse, err := http.Get(fmt.Sprintf("%s/%s/-/%s-%s.tgz", c.url, fullName, name, version))
	if err != nil {
		return nil, fmt.Errorf("error fetching tarball: %w", err)
	}
	defer fetchResponse.Body.Close()

	if fetchResponse.StatusCode != 200 {
		return nil, fmt.Errorf("registry error: %d", fetchResponse.StatusCode)
	}

	gzr, err := gzip.NewReader(fetchResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	res := GetPackageResponse{Files: make([]GetPackageResponseFile, 0)}
	for header, err := tr.Next(); err != io.EOF; header, err = tr.Next() {
		if err != nil {
			return nil, fmt.Errorf("failed to read tarball: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		trimmedPath := strings.TrimPrefix(header.Name, "package/")
		buf, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to read tarball: %w", err)
		}

		f := GetPackageResponseFile{
			Stream: buf,
			Path:   trimmedPath,
		}
		res.Files = append(res.Files, f)
	}

	return &res, nil
}
