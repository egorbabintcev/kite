package registry

type GetPackageResponseFile struct {
	Stream []byte
	Path   string
}

type GetPackageResponse struct {
	Files []GetPackageResponseFile
}
