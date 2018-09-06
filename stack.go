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
	t := struct {
		Name      string        `json:",omitempty"`
		Component *ComponentRef `json:",omitempty"`
		On        []string      `json:",omitempty"`
		Hooks     *StackHook    `json:",omitempty"`
	}{
		Name:      r.Name,
		Component: &r.Component,
	}
	for _, k := range r.On.nodeSets {
		t.On = append(t.On, k.Name)
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

type StackHook struct {
	Deploy   Hook
	Undeploy Hook
}

func (r StackHook) HasTasks() bool {
	return r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks()
}

func (r StackHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Deploy   *Hook `json:",omitempty"`
		Undeploy *Hook `json:",omitempty"`
	}{}

	if r.Deploy.HasTasks() {
		t.Deploy = &r.Deploy
	}
	if r.Undeploy.HasTasks() {
		t.Undeploy = &r.Undeploy
	}
	return json.Marshal(t)
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
