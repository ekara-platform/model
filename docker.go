package model

type Docker struct {
	docker map[string]interface{}
}

func createDocker(vErrs *ValidationErrors, p map[string]interface{}) Docker {
	ret := Docker{map[string]interface{}{}}
	for k, v := range p {
		ret.docker[k] = v
	}
	return ret
}

func (d Docker) AsMap() map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range d.docker {
		ret[k] = v
	}
	return ret
}

func (d *Docker) add(name string, value interface{}) {
	d.docker[name] = value
}
