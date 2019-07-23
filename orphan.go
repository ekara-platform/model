package model

type (
	// Orphans represents and indirect, unresolved, references used into
	// an environment descriptor
	Orphans struct {
		// Set used to avoid duplicated entries
		Refs map[string]struct{}
	}
)

func CreateOrphans() *Orphans {
	res := &Orphans{
		Refs: make(map[string]struct{}),
	}
	return res
}

func (or *Orphans) add(ref, kind string) {
	if ref != "" && kind != "" {
		or.Refs[ref+"-"+kind] = struct{}{}
	}
}

func (or *Orphans) AddAll(ors Orphans) {
	for key := range ors.Refs {
		if key != "" {
			or.Refs[key] = struct{}{}
		}
	}
}

func (or *Orphans) NoMoreAnOrhpan(ref string) {
	delete(or.Refs, ref)
}
