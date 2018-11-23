package model

import (
	"encoding/json"
	"errors"
)

type (
	//Stack represent an Stack installable on the built environment
	Stack struct {
		// The name of the stack
		Name string
		// The component containing the stack
		Component componentRef
		// The hooks linked to the stack lifecycle events
		Hooks StackHook
	}

	//Stacks represent all the stacks of an environment
	Stacks map[string]Stack
)

//DescType returns the Describable type of the stack
//  Hardcoded to : "Stack"
func (r Stack) DescType() string {
	return "Stack"
}

//DescName returns the Describable name of the stack
func (r Stack) DescName() string {
	return r.Name
}

// MarshalJSON returns the serialized content of the stack as JSON
func (r Stack) MarshalJSON() ([]byte, error) {
	t := struct {
		Name      string        `json:",omitempty"`
		Component *componentRef `json:",omitempty"`
		On        []string      `json:",omitempty"`
		Hooks     *StackHook    `json:",omitempty"`
	}{
		Name:      r.Name,
		Component: &r.Component,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

func (r Stack) validate() ValidationErrors {
	return ErrorOnInvalid(r.Component, r.Hooks)
}

func (r *Stack) merge(other Stack) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Component.merge(other.Component); err != nil {
		return err
	}
	return nil
}

func createStacks(env *Environment, yamlEnv *yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		res[name] = Stack{
			Name:      name,
			Component: createComponentRef(env, env.location.appendPath("stacks."+name+".component"), yamlStack.Component, false),
			Hooks: StackHook{
				Deploy:   createHook(env, env.location.appendPath("stacks."+name+".hooks.deploy"), yamlStack.Hooks.Deploy),
				Undeploy: createHook(env, env.location.appendPath("stacks."+name+".hooks.undeploy"), yamlStack.Hooks.Undeploy)}}
	}
	return res
}

func (r Stacks) merge(env *Environment, others Stacks) error {
	for id, s := range others {
		if stack, ok := r[id]; ok {
			if err := stack.merge(s); err != nil {
				return err
			}
		} else {
			s.Component.env = env
			r[id] = s
		}
	}
	return nil
}
