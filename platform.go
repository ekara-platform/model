package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base                       Base
	Parent                     Parent
	Components                 map[string]Component
	sortedDiscoveredComponents []ComponentLeaf
	SortedFetchedComponents    []string
	cRefs                      []ComponentReferencer
}

//ComponentLeaf keeps a reference between a component and the one
// declaring it
type ComponentLeaf struct {
	parentID  string
	Component Component
}

func createPlatform(yamlEnv *yamlEnvironment) (*Platform, error) {
	p := &Platform{}
	// Compute the component base for the environment
	base, e := CreateComponentBase(yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the base component : " + e.Error())
	}
	p.Base = base

	// Create the parent component
	parent, e := CreateParent(base, yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the parent : " + e.Error())
	}
	p.Parent = parent

	// Create other components of the environment
	p.sortedDiscoveredComponents = make([]ComponentLeaf, 0, 0)
	p.SortedFetchedComponents = make([]string, 0, 0)
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
	p.cRefs = make([]ComponentReferencer, 0, 0)
	return p, nil
}

//RegisterComponent register a new component under its parent ID
//If a component has no parent, like a main descriptor, then the ID should be ""
func (p *Platform) RegisterComponent(parent string, c Component) bool {
	res := false
	if _, ok := p.Components[c.Id]; ok {

		for i, v := range p.sortedDiscoveredComponents {
			if v.Component.Id == c.Id {
				p.sortedDiscoveredComponents = append(p.sortedDiscoveredComponents[:i], p.sortedDiscoveredComponents[i+1:]...)
				break
			}
		}

		res = true
	}
	p.sortedDiscoveredComponents = append(p.sortedDiscoveredComponents, ComponentLeaf{
		parentID:  parent,
		Component: c,
	})
	p.Components[c.Id] = c
	return res
}

//ToFetch provides a channel allowing to get the sorted components to fetch
func (p *Platform) ToFetch() (<-chan Component, int) {
	sD := p.sortedDiscoveredComponents
	ret := make(chan Component, len(sD))

	go func() {
		work := make([]ComponentLeaf, len(sD))
		copy(work, sD)
		lastDone := make([]string, 0, 0)

		// First we check the parent
		for i, n := range work {
			if n.Component.Id == EkaraComponentId {
				lastDone = append(lastDone, n.Component.Id)
				work = append(work[:i], work[i+1:]...)
				ret <- n.Component
				continue
			}
		}
		for len(work) > 0 {
			lastDoneCHildren := false
			for i, n := range work {
				if len(lastDone) > 0 {
					if n.parentID == lastDone[len(lastDone)-1] {
						ret <- n.Component
						work = append(work[:i], work[i+1:]...)
						lastDone = append(lastDone, n.Component.Id)
						lastDoneCHildren = true
						break
					}
				} else {
					ret <- n.Component
					work = append(work[:i], work[i+1:]...)
					lastDone = append(lastDone, n.Component.Id)
					break
				}
			}
			if !lastDoneCHildren {
				if len(lastDone) > 0 {
					lastDone = lastDone[:len(lastDone)-1]
				}
			}
		}
		close(ret)
	}()
	return ret, len(sD)
}

func (p *Platform) tagUsedComponent(cr ComponentReferencer) {
	for _, v := range p.cRefs {
		if cr.ComponentName() == v.ComponentName() {
			return
		}
	}
	p.cRefs = append(p.cRefs, cr)
}

// Used returns true if the component with the given Id is used into the environment.
func (p *Platform) Used(id string) bool {
	for _, cr := range p.cRefs {
		if cr.ComponentName() == id {
			return true
		}
	}
	return false
}

func (p Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range p.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

func (p *Platform) merge(other Platform) error {
	// We don't need to merge other.Components because they will be processed as
	// components to be registered by the fetching process into the engine
	for _, c := range other.cRefs {
		p.tagUsedComponent(c)
	}

	if p.Parent.Repository.Url == nil && other.Parent.Repository.Url != nil {
		p.Parent = other.Parent
	}
	return nil
}
