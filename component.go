package model

type Component struct {
	Repository string
	Version    Version
}

func createComponent(vErrs *ValidationErrors, repository string, version Version) Component {
	return Component{Repository: repository, Version: version}
}
