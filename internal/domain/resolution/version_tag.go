package resolution

const LatestTagName = "latest"

type VersionTag struct {
	name string
}

func NewVersionTag(name string) (VersionTag, error) {
	if name == "" {
		return VersionTag{}, ErrVersionTagEmpty
	}

	return VersionTag{name: name}, nil
}

func (vt VersionTag) IsLatest() bool {
	return vt.name == LatestTagName
}
