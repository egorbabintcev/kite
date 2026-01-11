package resolution

import "fmt"

const LatestTagName = "latest"

type VersionTag struct {
	name string
}

func NewVersionTag(name string) (VersionTag, error) {
	if name == "" {
		return VersionTag{}, fmt.Errorf("tag name cannot be empty")
	}

	return VersionTag{name: name}, nil
}

func (vt VersionTag) IsLatestTag() bool {
	return vt.name == LatestTagName
}
