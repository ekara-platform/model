package model

import (
	"sort"
)

type (

	//ReferencedComponent keeps a reference between a component delacation and the one
	// declaring it
	ReferencedComponent struct {
		Owner     string
		Component Component
	}

	// ReferencedComponents represents a manager of declarations of components
	// into an environment descriptor
	ReferencedComponents struct {
		Refs []ReferencedComponent
	}
)

// CreateReferencedComponents return an initialized manager
func CreateReferencedComponents() *ReferencedComponents {
	res := &ReferencedComponents{
		Refs: make([]ReferencedComponent, 0, 0),
	}
	return res
}

// Add a new component reference.
// It will return false is a component with the same id or the same repository
// as already been registered.
func (rc *ReferencedComponents) add(owner string, c Component) bool {
	for _, ve := range rc.Refs {
		// If the component has been already referenced then it won't be added/overwritten
		if ve.Component.Repository.Url.String() == c.Repository.Url.String() || ve.Component.Id == c.Id {
			return false
		}
	}
	rc.Refs = append(rc.Refs, ReferencedComponent{
		Owner:     owner,
		Component: c,
	})
	return true
}

// AddReference adds a new referenced component.
// It will return false is a component with the same id or the same repository
// as already been registered.
func (rc *ReferencedComponents) AddReference(ref ReferencedComponent) bool {
	return rc.add(ref.Owner, ref.Component)
}

//IdReferenced return true if a component with the given id is referenced
func (rc *ReferencedComponents) IdReferenced(id string) bool {
	for _, ve := range rc.Refs {
		if ve.Component.Id == id {
			return true
		}
	}
	return false
}

//Clean cleans the referenced component base on the list of used ones
func (rc *ReferencedComponents) Clean(used UsedReferences) {
	var cleaned []ReferencedComponent
	for _, v := range rc.Refs {
		for k := range used.Refs {
			if v.Component.Id == k {
				cleaned = append(cleaned, v)
				break
			}
		}
	}
	rc.Refs = cleaned
}

//Sorted Returns the referenced components sorted in alphabetical order
// based on their names
func (rc *ReferencedComponents) Sorted() []Component {
	var res []Component
	if len(rc.Refs) > 0 {
		var keys []string
		kVs := make(map[string]Component)
		for _, v := range rc.Refs {
			keys = append(keys, v.Component.Id)
			kVs[v.Component.Id] = v.Component
		}
		sort.Strings(keys)

		for _, k := range keys {
			if val, ok := kVs[k]; ok {
				res = append(res, val)
			}
		}
	}
	return res
}
