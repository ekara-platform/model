package model

type Stack struct {
	// The environment holding the stack
	root *Environment
	// The labels associated to the stack
	Labels
	// The repository/version of the stack
	Component

	// The name of the stack
	Name string
	// The specifications on where the stack is supposed to deployed
	DeployOn NodeSetRef

	// The hooks linked to the stack lifecycle
	Hooks struct {
		Deploy   Hook
		Undeploy Hook
	}
}

func createStacks(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Stack {
	res := map[string]Stack{}
	if yamlEnv.Stacks == nil || len(yamlEnv.Stacks) == 0 {
		vErrs.AddWarning("no stack specified", "stacks")
	} else {
		for name, yamlStack := range yamlEnv.Stacks {
			stack := Stack{
				root:   env,
				Labels: createLabels(vErrs, yamlStack.Labels...),
				Name:   name}

			stack.Component = createComponent(vErrs, env.Lagoon, "stacks."+name, yamlStack.Repository, yamlStack.Version)
			stack.DeployOn = createNodeSetRef(vErrs, env, "stacks."+name+".deployOn", yamlStack.DeployOn...)
			stack.Hooks.Deploy = createHook(vErrs, env.Tasks, "stacks."+name+".hooks.deploy", yamlStack.Hooks.Deploy)
			stack.Hooks.Undeploy = createHook(vErrs, env.Tasks, "stacks."+name+".hooks.undeploy", yamlStack.Hooks.Undeploy)

			res[name] = stack
		}
	}
	return res
}
