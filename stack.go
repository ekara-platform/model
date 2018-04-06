package descriptor

type Stack struct {
	root *Environment
	Labels
	Component

	Name string

	Hooks struct {
		Deploy   Hook
		Undeploy Hook
	}
}

func createStacks(env *Environment, yamlEnv *yamlEnvironment) (res map[string]Stack, err error) {
	res = map[string]Stack{}
	for name, yamlStack := range yamlEnv.Stacks {
		stack := Stack{
			root:   env,
			Labels: createLabels(yamlStack.Labels...),
			Name:   name}

		stack.Component, err = createComponent(yamlStack.Repository, yamlStack.Version)
		if err != nil {
			return
		}
		stack.Hooks.Deploy, err = createHook(env.Tasks, yamlStack.Hooks.Deploy)
		if err != nil {
			return
		}
		stack.Hooks.Undeploy, err = createHook(env.Tasks, yamlStack.Hooks.Undeploy)
		if err != nil {
			return
		}

		res[name] = stack
	}
	return
}
