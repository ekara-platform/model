package model

type Parameters map[string]interface{}

func createParameters(src map[string]interface{}) Parameters {
	dst := make(map[string]interface{})
	mergeParams(dst, src)
	return dst
}

func (p Parameters) inherit(parents ...map[string]interface{}) Parameters {
	dst := make(map[string]interface{})
	for i := len(parents) - 1; i >= 0; i-- {
		mergeParams(dst, parents[i])
	}
	mergeParams(dst, p)
	return dst
}

func mergeParams(dst map[string]interface{}, src map[string]interface{}) {
	// TODO take into account map merging
	for k, v := range src {
		dst[k] = v
	}
}
