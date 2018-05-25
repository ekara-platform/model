package model

type attributes map[string]interface{}

func createAttributes(p attributes, herited attributes) attributes {
	if herited != nil {
		return mergeMap(p, herited)
	}
	m := map[string]interface{}{}
	for k, v := range p {
		m[k] = v
	}
	return m
}

func mergeMap(m map[string]interface{}, herited map[string]interface{}) attributes {
	r := make(map[string]interface{})
	if herited != nil {
		for k, v := range herited {
			r[k] = v
		}
	}

	for k, v := range m {
		r[k] = v
	}
	return r
}

func (m attributes) copy() attributes {
	r := map[string]interface{}{}
	for k, v := range m {
		r[k] = v
	}
	return r
}
