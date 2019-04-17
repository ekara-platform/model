package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base         Base
	Distribution Component
	Components   map[string]Component
}

func createPlatform(yamlEnv *yamlEnvironment) (Platform, error) {
	components := map[string]Component{}
	p := Platform{}
	// Compute the component base for the environment
	base, e := CreateComponentBase(yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the base component : " + e.Error())
	}

	// Create the distribution component (mandatory)
	ekaraRepo := yamlEnv.Ekara.Distribution.Repository
	if ekaraRepo == "" {
		ekaraRepo = EkaraComponentRepo
	}
	repoDist, e := CreateRepository(base, ekaraRepo, yamlEnv.Ekara.Distribution.Ref, "")
	if e != nil {
		return p, errors.New("invalid distribution: " + e.Error())
	}
	ekaraComponent := CreateComponent(EkaraComponentId, repoDist)

	setCredentials(&ekaraComponent, yamlEnv.Ekara.Distribution)

	// Create other components of the environment
	for name, yamlC := range yamlEnv.Ekara.Components {
		repo, e := CreateRepository(base, yamlC.Repository, yamlC.Ref, "")
		if e != nil {
			return p, errors.New("Error creating the repository: " + e.Error())
		}
		component := CreateComponent(name, repo)
		setCredentials(&component, yamlC)
		components[name] = component
	}

	p.Base = base
	p.Distribution = ekaraComponent
	p.Components = components
	return p, nil
}

func (r Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range r.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

func (r *Platform) merge(other Platform) error {
	for id, c := range other.Components {
		if _, ok := r.Components[id]; !ok {
			r.Components[id] = c
		}
	}
	return nil
}

func setCredentials(component *Component, yamlComponent yamlComponent) {
	if len(yamlComponent.Auth) > 0 {
		component.Authentication = createParameters(yamlComponent.Auth)
	}
}
