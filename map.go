package model

type attributes map[string]interface{}
type envvars map[string]string

func createAttributes(p attributes, herited attributes) attributes {
	if herited != nil {
		return mergeAttributeMap(p, herited)
	}
	m := map[string]interface{}{}
	for k, v := range p {
		m[k] = v
	}
	return m
}

func createEnvvars(p envvars, herited envvars) envvars {
	if herited != nil {
		return mergeEnvvarsMap(p, herited)
	}
	m := map[string]string{}
	for k, v := range p {
		m[k] = v
	}
	return m
}

func mergeAttributeMap(m map[string]interface{}, herited map[string]interface{}) attributes {
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

func mergeEnvvarsMap(m map[string]string, herited map[string]string) envvars {
	r := make(map[string]string)
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

func (m envvars) copy() envvars {
	r := map[string]string{}
	for k, v := range m {
		r[k] = v
	}
	return r
}
