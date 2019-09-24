package model

type (

	// EnvironmentReferences represents a light Ekara environment, used to unmarshal component references only
	EnvironmentReferences struct {
		Ekara    yamlEkara
		yamlVars `yaml:",inline"`

		OrchestratorRefs struct {
			Component string
		} `yaml:"orchestrator"`

		ProvidersRefs map[string]struct {
			Component string
		} `yaml:"providers"`

		NodesRefs map[string]struct {
			Provider struct {
				Component string `yaml:"name"`
			}
		} `yaml:"nodes"`

		StacksRefs map[string]struct {
			Component string
		} `yaml:"stacks"`

		TasksRefs map[string]struct {
			Component string
		} `yaml:"tasks"`
	}
)

// Uses returns the references of the components used into the environment
func (er EnvironmentReferences) Uses(previousO *Orphans) (*UsedReferences, *Orphans) {
	res := CreateUsedReferences()
	orphans := CreateOrphans()

	res.add(er.OrchestratorRefs.Component)

	for k := range previousO.Refs {
		key, kind := previousO.KeyType(k)
		if kind == "provider" {
			for pKey, pval := range er.ProvidersRefs {
				if key == pKey {
					res.add(pval.Component)
					previousO.NoMoreAnOrhpan(k)
					break
				}
			}
		}
	}

	// We pass through the node sets to get the used providers
	for _, val := range er.NodesRefs {
		located := false
		for key, pval := range er.ProvidersRefs {
			// An Orphan, of type provider, must be added
			// if the provider component is unknown ( for example
			// when the provider in overwritten), of if the component
			// is known but not yet referenced
			if val.Provider.Component == key {
				// component unknown
				if pval.Component == "" {
					orphans.new(key, "provider")
				} else {
					res.add(pval.Component)
				}
				located = true
			}
		}
		// component not yet referenced
		if !located {
			orphans.new(val.Provider.Component, "provider")
		}
	}
	for _, val := range er.StacksRefs {
		res.add(val.Component)
	}
	for _, val := range er.TasksRefs {
		res.add(val.Component)
	}
	return res, orphans
}

// References returns the components referenced into the environment
func (er EnvironmentReferences) References(owner string) (*ReferencedComponents, error) {
	res := CreateReferencedComponents()

	ekara, err := createPlatform(er.Ekara)
	if err != nil {
		return res, err
	}

	for _, val := range ekara.Components {
		res.add(owner, val)
	}
	return res, nil
}

// Parent returns the parent of the component
func (er EnvironmentReferences) Parent() (parent Parent, err error) {
	var parentBase Base
	parentBase, err = CreateBase(er.Ekara.Base)
	if err != nil {
		return
	}
	parent, err = CreateParent(parentBase, er.Ekara)
	return
}
