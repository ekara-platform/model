package model

import (
	"errors"
)

//Platform the platform used to build an environment
type Platform struct {
	Base                       Base
	Distribution               Distribution
	Components                 map[string]Component
	sortedDiscoveredComponents []ComponentLeaf
	SortedFetchedComponents    []string
	cRefs                      []ComponentReferencer
}

type ComponentLeaf struct {
	ParentId  string
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

	// Create the distribution component (mandatory)
	dist, e := CreateDistribution(base, yamlEnv)
	if e != nil {
		return p, errors.New("Error creating the distribution : " + e.Error())
	}
	p.Distribution = dist

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

func (p *Platform) RegisterComponent(parent string, c Component) {
	p.sortedDiscoveredComponents = append(p.sortedDiscoveredComponents, ComponentLeaf{
		ParentId:  parent,
		Component: c,
	})

	if _, ok := p.Components[c.Id]; !ok {
		p.Components[c.Id] = c
	}
}

func (p *Platform) ToFetch() (<-chan Component, int) {
	sD := p.sortedDiscoveredComponents
	ret := make(chan Component, len(sD))

	go func() {
		work := make([]ComponentLeaf, len(sD))
		copy(work, sD)
		lastDone := make([]string, 0, 0)

		// First we check the distribution
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
					if n.ParentId == lastDone[len(lastDone)-1] {
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
	p.cRefs = append(p.cRefs, cr)
}

// UsedComponents returns an array of components effectively in used throughout the descriptor.
// TODO SUPPOSE TO BE DELETED
func (p *Platform) UsedComponents() ([]Component, error) {
	res := make([]Component, 0, 0)
	temp := make(map[string]Component)
	for _, cr := range p.cRefs {
		name := cr.ComponentName()
		if name != "" {
			c, err := cr.Component()
			if err != nil {
				continue
			}
			temp[name] = c
		}
	}
	for _, c := range temp {
		res = append(res, c)
	}
	return res, nil
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
	for id, c := range other.Components {
		if id != "" {
			if _, ok := p.Components[id]; !ok {
				p.Components[id] = c
			}
		}
	}

	for _, c := range other.cRefs {
		p.tagUsedComponent(c)
	}

	if p.Distribution.Repository.Url == nil && other.Distribution.Repository.Url != nil {
		p.Distribution = other.Distribution
	}
	return nil
}
