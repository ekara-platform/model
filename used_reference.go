package model

type (
	// UsedReferences represents a manager of references to components used into
	// an environment descriptor
	UsedReferences struct {
		// Set used to avoid duplicated entries
		Refs map[string]struct{}
	}
)

//CreateUsedReferences returns an initialized manager
func CreateUsedReferences() *UsedReferences {
	res := &UsedReferences{
		Refs: make(map[string]struct{}),
	}
	return res
}

// add a new reference on a used component, if a reference with the same
// id has already been registered then it will be ignored
func (ur *UsedReferences) add(id string) {
	if id != "" {
		ur.Refs[id] = struct{}{}
	}
}

// AddAll adds a new references on used components, if a reference with the same
// id has already been registered then it will be ignored
func (ur *UsedReferences) AddAll(ids UsedReferences) {
	for key := range ids.Refs {
		if key != "" {
			ur.Refs[key] = struct{}{}
		}
	}
}

// IdUsed returns true if a component with the provided id has been referenced as used
func (ur *UsedReferences) IdUsed(id string) bool {
	_, ok := ur.Refs[id]
	return ok
}
