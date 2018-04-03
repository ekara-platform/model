package descriptor

type NodeSet struct {
	root *Environment
	Labels

	Name      string
	Provider  ProviderRef
	Instances int

	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
}

func createNodeSets(env *Environment, yamlEnv *yamlEnvironment) (res map[string]NodeSet, err error) {
	res = map[string]NodeSet{}
	for name, yamlNodeSet := range yamlEnv.Nodes {
		nodeSet := NodeSet{
			root:      env,
			Labels:    createLabels(yamlNodeSet.Labels...),
			Name:      name,
			Instances: yamlNodeSet.Instances}

		nodeSet.Provider, err = createProviderRef(env, yamlNodeSet.Provider)
		if err != nil {
			return
		}

		nodeSet.Hooks.Provision, err = createHook(env.Tasks, yamlNodeSet.Hooks.Provision)
		if err != nil {
			return
		}
		nodeSet.Hooks.Destroy, err = createHook(env.Tasks, yamlNodeSet.Hooks.Destroy)
		if err != nil {
			return
		}

		res[name] = nodeSet
	}
	return
}
