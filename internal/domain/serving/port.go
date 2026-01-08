package serving

type PackageFileContentProvider interface {
	Get(path string) (PackageFileContent, error)
}
