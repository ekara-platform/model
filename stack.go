package model

import (
	"encoding/json"
)

type Stack struct {
	// The name of the stack
	Name string
	// The component containing the stack
	Component ComponentRef
	// The node sets where the stack should be deployed
	On NodeSetRef
	// The hooks linked to the stack lifecycle
	Hooks StackHook
}

func (r Stack) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string        `json:",omitempty"`
		Component *ComponentRef `json:",omitempty"`
		On        *NodeSetRef   `json:",omitempty"`
		Hooks     *StackHook    `json:",omitempty"`
	}{
		Name:      r.Name,
		Component: &r.Component,
		On:        &r.On,
		Hooks:     &r.Hooks,
	})
}

type StackHook struct {
	Deploy   Hook
	Undeploy Hook
}

func (r StackHook) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Deploy   *Hook `json:",omitempty"`
		Undeploy *Hook `json:",omitempty"`
	}{
		Deploy:   &r.Deploy,
		Undeploy: &r.Undeploy,
	})
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
				Hooks: StackHook{
					Deploy:   createHook(vErrs, "stacks."+name+".hooks.deploy", env, yamlStack.Hooks.Deploy),
					Undeploy: createHook(vErrs, "stacks."+name+".hooks.undeploy", env, yamlStack.Hooks.Undeploy)}}
		}
	}
	return res
}
