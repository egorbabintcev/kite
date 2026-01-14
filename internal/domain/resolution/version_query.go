package resolution

type VersionQuery struct {
	Range *VersionRange
	Tag   *VersionTag
}

func NewVersionQuery(raw string) (VersionQuery, error) {
	if raw == "" {
		t, err := NewVersionTag(LatestTagName)
		if err != nil {
			return VersionQuery{}, err
		}

		return VersionQuery{Tag: &t}, nil
	}

	r, err := NewVersionRange(raw)
	if err != nil {
		t, err := NewVersionTag(raw)
		if err != nil {
			return VersionQuery{}, ErrVersionQueryInvalid
		}

		return VersionQuery{Tag: &t}, nil
	}

	return VersionQuery{Range: &r}, nil
}

func (vq VersionQuery) Equal(other VersionQuery) bool {
	if vq.Range != nil && other.Range != nil {
		return vq.Range.Equal(*other.Range)
	}

	if vq.Tag != nil && other.Tag != nil {
		return *vq.Tag == *other.Tag
	}

	return false
}
