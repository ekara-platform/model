package model

type (
	ComponentRef struct {
		ref       string
		mandatory bool
		env       *Environment
		location  DescriptorLocation
	}
)

func (r ComponentRef) Reference() Reference {
	result := make(map[string]interface{})
	for k, v := range r.env.Ekara.Components {
		result[k] = v
	}
	return Reference{
		Id:        r.ref,
		Type:      "component",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

func (r *ComponentRef) merge(other ComponentRef) error {
	if r.ref == "" {
		r.ref = other.ref
	}
	return nil
}

func (r ComponentRef) Resolve() (Component, error) {
	validationErrors := ErrorOn(r)
	if validationErrors.HasErrors() {
		return Component{}, validationErrors
	}
	return r.env.Ekara.Components[r.ref], nil
}

func createComponentRef(env *Environment, location DescriptorLocation, componentRef string, mandatory bool) ComponentRef {
	return ComponentRef{
		env:       env,
		location:  location,
		ref:       componentRef,
		mandatory: mandatory,
	}
}
