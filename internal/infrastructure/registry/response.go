package registry

type GetPackageResponseFile struct {
	Content []byte
	Path    string
}

type GetPackageResponse struct {
	Files []GetPackageResponseFile
}

type GetMetadataResponseMetadata struct {
	Versions []string
	Tags     map[string]string
}

type GetMetadataResponse struct {
	Metadata GetMetadataResponseMetadata
}
