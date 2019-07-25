package model

import "strings"

type (
	// Orphans represents and indirect, unresolved, references used into
	// an environment descriptor
	Orphans struct {
		// Set used to avoid duplicated entries
		Refs map[string]struct{}
	}
)

// CreateOrphans return an initialized manager
func CreateOrphans() *Orphans {
	res := &Orphans{
		Refs: make(map[string]struct{}),
	}
	return res
}

//new adds a new orphan.
// The new added orphan will be formated like this : ref-king
func (or *Orphans) new(ref, kind string) {
	if ref != "" && kind != "" {
		or.Refs[ref+"-"+kind] = struct{}{}
	}
}

//AddReference adds a new reference to orphan which has already been formated like : ref-king
func (or *Orphans) AddReference(id string) {
	if id != "" {
		or.Refs[id] = struct{}{}
	}
}

//NoMoreAnOrhpan remove an orphan, which has already been formated like : ref-king.
func (or *Orphans) NoMoreAnOrhpan(ref string) {
	delete(or.Refs, ref)
}

//KeyType returns the key and the type of an.
func (or *Orphans) KeyType(s string) (key string, kind string) {
	i := strings.LastIndex(s, "-")
	key = s[0:i]
	kind = s[i+1:]
	return
}
