package model

type EnvVars map[string]string

func createEnvVars(src map[string]string) EnvVars {
	dst := make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (r EnvVars) inherits(parent map[string]string) EnvVars {
	dst := make(map[string]string)
	for k, v := range parent {
		dst[k] = v
	}
	for k, v := range r {
		dst[k] = v
	}
	return dst
}
