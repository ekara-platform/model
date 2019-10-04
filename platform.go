package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base       Base
	Parent     Parent
	HasParent  bool
	Components map[string]Component
}

func createPlatform(yamlEkara yamlEkara) (Platform, error) {
	p := Platform{}
	// Compute the component base for the environment
	base, e := CreateComponentBase(yamlEkara)
	if e != nil {
		return p, errors.New("Error creating the base component : " + e.Error())
	}
	p.Base = base

	// Create the parent component
	parent, hasParent, e := CreateParent(base, yamlEkara)
	if e != nil {
		return p, errors.New("Error creating the parent : " + e.Error())
	}
	p.HasParent = hasParent
	p.Parent = parent

	// Create other components of the environment
	components := map[string]Component{}
	for name, yamlC := range yamlEkara.Components {
		repo, e := CreateRepository(base, yamlC.Repository, yamlC.Ref, "")
		if e != nil {
			return p, errors.New("Error creating the repository: " + e.Error())
		}
		repo.setAuthentication(yamlC)
		components[name] = CreateComponent(name, repo)
	}

	p.Components = components
	return p, nil
}

func (p Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range p.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

//KeepTemplates Stores the template into the given component
func (p Platform) KeepTemplates(c Component, templates Patterns) {
	if len(templates.Content) > 0 {
		comp := p.Components[c.Id]
		comp.Templates = templates
		p.Components[c.Id] = comp
	}
}

//AddComponent Add the given component to the platform
func (p *Platform) AddComponent(c Component) {
	p.Components[c.Id] = c
}
