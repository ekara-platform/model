package model

type Labels map[string]string

func (r Labels) inherits(parent Labels) Labels {
	dst := make(map[string]string)
	for k, v := range parent {
		dst[k] = v
	}
	for k, v := range r {
		dst[k] = v
	}
	return dst
}
