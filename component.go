package descriptor

type Component struct {
	Repository string
	Version    Version
}

func createComponent(repository string, version string) (res Component, err error) {
	res = Component{Repository: repository}
	res.Version, err = createVersion(version)
	return
}
