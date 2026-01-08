package shared_kernel

import "fmt"

type PackageID struct {
	scope string
	name  string
}

func NewPackageID(scope string, name string) (*PackageID, error) {
	if len(scope) < 2 {
		return nil, fmt.Errorf("scope must be at least 2 characters")
	}

	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	return &PackageID{
		scope: scope,
		name:  name,
	}, nil
}

func (p *PackageID) String() string {
	return fmt.Sprintf("@%s/%s", p.scope, p.name)
}
