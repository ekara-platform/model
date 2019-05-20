package model

//EnvVars Represents environment variable
type EnvVars map[string]string

func createEnvVars(src map[string]string) (EnvVars, error) {
	dst := make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst, nil
}

func (r EnvVars) inherit(parent map[string]string) (EnvVars, error) {
	dst := make(map[string]string)
	for k, v := range parent {
		dst[k] = v
	}
	for k, v := range r {
		dst[k] = v
	}
	return dst, nil
}
