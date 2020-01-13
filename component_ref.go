package model

import (
	"fmt"
)

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
)

func createComponentRef(env *Environment, location DescriptorLocation, ref string, mandatory bool) componentRef {
	return componentRef{
		env:       env,
		location:  location,
		ref:       ref,
		mandatory: mandatory,
	}
}

func (r *componentRef) customize(with componentRef) error {
	if with.ref != "" {
		r.ref = with.ref
	}
	r.mandatory = with.mandatory
	if !with.location.empty() {
		r.location = with.location
	}

	return nil
}

func (r componentRef) resolve() (Component, error) {
	if val, ok := r.env.ekara.Components[r.ref]; ok {
		return val, nil
	}
	return Component{}, fmt.Errorf(unknownComponentRefError, r.ref)
}

func (r componentRef) validationDetails() refValidationDetails {
	result := make(map[string]interface{})
	for k, v := range r.env.ekara.Components {
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
