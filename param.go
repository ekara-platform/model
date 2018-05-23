package model

type Parameters struct {
	parameters map[string]interface{}
}

func createParameters(vErrs *ValidationErrors, p map[string]interface{}) Parameters {
	ret := Parameters{map[string]interface{}{}}
	for k, v := range p {
		ret.parameters[k] = v
	}
	return ret
}

func (p Parameters) AsMap() map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range p.parameters {
		ret[k] = v
	}
	return ret
}

func (p *Parameters) add(name string, value interface{}) {
	p.parameters[name] = value
}
