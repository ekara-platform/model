package model

type Stack struct {
	// The name of the stack
	Name string
	// The component containing the stack
	Component ComponentRef
	// The node sets where the stack should be deployed
	On NodeSetRef
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
			res[name] = Stack{
				Name:      name,
				Component: createComponentRef(vErrs, env.Lagoon.Components, "stacks."+name+".component", yamlStack.Component),
				On:        createNodeSetRef(vErrs, env, "stacks."+name+".on", yamlStack.On...),
				Hooks: struct {
					Deploy   Hook
					Undeploy Hook
				}{
					Deploy:   createHook(vErrs, "stacks."+name+".hooks.deploy", env, yamlStack.Hooks.Deploy),
					Undeploy: createHook(vErrs, "stacks."+name+".hooks.undeploy", env, yamlStack.Hooks.Undeploy)}}
		}
	}
	return res
}
