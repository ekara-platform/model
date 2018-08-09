package model

type EnvVars map[string]string

func createEnvVars(src map[string]string) EnvVars {
	dst := make(map[string]string)
	mergeEnvVars(dst, src)
	return dst
}

func (e EnvVars) inherit(parents ...map[string]string) EnvVars {
	dst := make(map[string]string)
	for i := len(parents) - 1; i >= 0; i-- {
		mergeEnvVars(dst, parents[i])
	}
	mergeEnvVars(dst, e)
	return dst
}

func mergeEnvVars(dst map[string]string, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}
