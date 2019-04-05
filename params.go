package model

//Parameters represents the parameters coming from a descriptor
type Parameters map[string]interface{}

func createParameters(src map[string]interface{}) Parameters {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (r Parameters) inherits(parent map[string]interface{}) Parameters {
	dst := make(map[string]interface{})
	for k, v := range parent {
		dst[k] = v
	}
	for k, v := range r {
		dst[k] = v
	}
	return dst
}
