package model

import "fmt"

const (
	unknownComponentRefError = "component cannot be resolved: %s"
)

type (
	//componentRef represents a reference to a component
	componentRef struct {
		//ref specifies id of the component
		ref string
		//mandatory indicates if the reference is mandatory
		mandatory bool
		//env specifies the environment holding the referenced component
		env *Environment
		//location indicates where the reference is located into the descriptor
		location DescriptorLocation
	}

	ComponentReferencer interface {
		Component() (Component, error)
		ComponentName() string
	}
)

func createComponentRef(env *Environment, location DescriptorLocation, ref string, mandatory bool) componentRef {
	return componentRef{
		env:       env,
		location:  location,
		ref:       ref,
		mandatory: mandatory,
	}
}

func (r *componentRef) merge(other componentRef) error {
	if r.ref == "" {
		r.ref = other.ref
	}
	return nil
}

func (r componentRef) resolve() (Component, error) {
	if val, ok := r.env.Ekara.Components[r.ref]; ok {
		return val, nil
	}
	return Component{}, fmt.Errorf(unknownComponentRefError, r.ref)
}

func (r componentRef) validationDetails() refValidationDetails {
	result := make(map[string]interface{})
	for k, v := range r.env.Ekara.Components {
		result[k] = v
	}
	return refValidationDetails{
		Id:        r.ref,
		Type:      "component",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}
