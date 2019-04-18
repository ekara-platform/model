package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base           Base
	Distribution   Distribution
	Components     map[string]Component
	UsedComponents []componentRef
}

func createPlatform(yamlEnv *yamlEnvironment) (*Platform, error) {

	p := &Platform{}
	// Compute the component base for the environment
	base, e := CreateComponentBase(yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the base component : " + e.Error())
	}
	p.Base = base

	// Create the distribution component (mandatory)
	dist, e := CreateDistribution(base, yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the distribution : " + e.Error())
	}
	p.Distribution = dist

	// Create other components of the environment
	components := map[string]Component{}
	for name, yamlC := range yamlEnv.Ekara.Components {
		repo, e := CreateRepository(base, yamlC.Repository, yamlC.Ref, "")
		if e != nil {
			return p, errors.New("Error creating the repository: " + e.Error())
		}
		repo.setAuthentication(yamlC)
		component := CreateComponent(name, repo)

		components[name] = component
	}

	p.Components = components
	p.UsedComponents = make([]componentRef, 0, 0)
	return p, nil
}

func (p *Platform) tagUsedComponent(cr componentRef) {
	if cr.ref != "" {
		for _, u := range p.UsedComponents {
			if u.ref == cr.ref {
				return
			}
		}
		p.UsedComponents = append(p.UsedComponents, cr)
	}
}

func (p Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range p.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

func (p *Platform) merge(other Platform) error {
	for id, c := range other.Components {
		if id != "" {
			if _, ok := p.Components[id]; !ok {
				p.Components[id] = c
			}
		}
	}

	for _, c := range other.UsedComponents {
		if c.ref != "" {
			p.tagUsedComponent(c)
		}
	}

	if p.Distribution.Repository.Url == nil && other.Distribution.Repository.Url != nil {
		p.Distribution = other.Distribution
	}
	return nil
}
