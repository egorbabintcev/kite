package registry

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
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
	fullName := c.resolveFullName(scope, name)

	url := fmt.Sprintf("%s/%s/-/%s-%s.tgz", c.url, fullName, name, version)
	fetchRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	fetchResponse, err := http.DefaultClient.Do(fetchRequest)
	if err != nil {
		return nil, fmt.Errorf("error fetching package archive: %w", err)
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
			Content: buf,
			Path:    trimmedPath,
		}
		res.Files = append(res.Files, f)
	}

	return &res, nil
}

func (c *HttpClient) FetchMetadata(ctx context.Context, scope, name string) (*GetMetadataResponse, error) {
	fullName := c.resolveFullName(scope, name)

	url := fmt.Sprintf("%s/%s", c.url, fullName)
	fetchRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	fetchResponse, err := http.DefaultClient.Do(fetchRequest)
	if err != nil {
		return nil, fmt.Errorf("error fetching package metadata: %w", err)
	}
	defer fetchResponse.Body.Close()

	var resData struct {
		Versions map[string]json.RawMessage `json:"versions"`
		Tags     map[string]string          `json:"dist-tags"`
	}
	err = json.NewDecoder(fetchResponse.Body).Decode(&resData)
	if err != nil {
		return nil, fmt.Errorf("error parsing package metadata: %w", err)
	}

	versions := make([]string, 0, len(resData.Versions))
	for v := range resData.Versions {
		versions = append(versions, v)
	}

	meta := GetMetadataResponseMetadata{
		Versions: versions,
		Tags:     resData.Tags,
	}

	return &GetMetadataResponse{
		Metadata: meta,
	}, nil
}

func (c *HttpClient) resolveFullName(scope, name string) string {
	if scope != "" {
		return fmt.Sprintf("@%s/%s", scope, name)
	}

	return name
}
