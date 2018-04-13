package model

type Parameters struct {
	parameters map[string]string
}

func createParameters(vErrs *ValidationErrors, p map[string]string) Parameters {
	ret := Parameters{map[string]string{}}
	for k, v := range p {
		ret.parameters[k] = v
	}
	return ret
}

func (p Parameters) AsMap() map[string]string {
	ret := map[string]string{}
	for k, v := range p.parameters {
		ret[k] = v
	}
	return ret
}

func (p *Parameters) add(name string, value string) {
	p.parameters[name] = value
}
