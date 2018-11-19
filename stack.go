package model

import (
	"encoding/json"
	"errors"
)

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

type Stack struct {
	// The name of the stack
	Name string
	// The component containing the stack
	Component ComponentRef
	// The hooks linked to the stack lifecycle
	Hooks StackHook
}

func (r Stack) DescType() string {
	return "Stack"
}

func (r Stack) DescName() string {
	return r.Name
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
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

func (r Stack) validate() ValidationErrors {
	vErrs := r.Component.validate()
	vErrs.merge(r.Hooks.Deploy.validate())
	vErrs.merge(r.Hooks.Undeploy.validate())
	return vErrs
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

type Stacks map[string]Stack

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

func (r Stacks) merge(env *Environment, other Stacks) error {
	for id, s := range other {
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
