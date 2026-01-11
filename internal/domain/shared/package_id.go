package shared

import "fmt"

type PackageID struct {
	scope string
	name  string
}

func NewPackageID(scope string, name string) (PackageID, error) {
	// Empty scope is allowed too
	if len(scope) == 1 {
		return PackageID{}, fmt.Errorf("scope must be at least 2 characters")
	}

	if name == "" {
		return PackageID{}, fmt.Errorf("name cannot be empty")
	}

	return PackageID{
		scope: scope,
		name:  name,
	}, nil
}

func (p PackageID) String() string {
	if p.scope != "" {
		return fmt.Sprintf("@%s/%s", p.scope, p.name)
	}
	return p.name
}
