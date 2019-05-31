package versionfilter

type Versions []*Version

func NewVersions(items []string) (Versions, error) {
	versions := Versions{}
	for _, s := range items {
		v, err := NewVersion(s)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (versions Versions) ToStrings() []string {
	strings := []string{}
	for _, v := range versions {
		strings = append(strings, v.Original)
	}
	return strings
}

func (versions Versions) Len() int {
	return len(versions)
}

func (versions Versions) Less(i, j int) bool {
	value, err := versions[i].Compare(versions[j])
	if err != nil {
		panic(err)
	}
	return value < 0
}

func (versions Versions) Swap(i, j int) {
	tmp := versions[i]
	versions[i] = versions[j]
	versions[j] = tmp
}
