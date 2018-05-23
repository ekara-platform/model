package model

type Stack struct {
	root *Environment
	Labels
	Component

	Name     string
	DeployOn NodeSetRef

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

			stack.Component = createComponent(vErrs, env, "stacks."+name, yamlStack.Repository, yamlStack.Version)
			stack.DeployOn = createNodeSetRef(vErrs, env, "stacks."+name+".deployOn", yamlStack.DeployOn...)
			stack.Hooks.Deploy = createHook(vErrs, env.Tasks, "stacks."+name+".hooks.deploy", yamlStack.Hooks.Deploy)
			stack.Hooks.Undeploy = createHook(vErrs, env.Tasks, "stacks."+name+".hooks.undeploy", yamlStack.Hooks.Undeploy)

			res[name] = stack
		}
	}
	return res
}
