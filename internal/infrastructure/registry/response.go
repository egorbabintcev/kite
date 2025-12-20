package registry

type GetPackageResponseFile struct {
	Content []byte
	Path    string
}

type GetPackageResponse struct {
	Files []GetPackageResponseFile
}
