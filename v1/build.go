package v1

// Build ...
type Build struct {
	sha1    string
	version string
}

// NewBuild ...
func NewBuild(sha1, version string) *Build {
	return &Build{
		sha1,
		version,
	}
}

func (b *Build) getSha1() string {
	return b.sha1
}

func (b *Build) getVersion() string {
	return b.version
}
