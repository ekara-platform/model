package model

import (
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
		Hooks      StackHook
		Parameters Parameters
		EnvVars    EnvVars
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
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)

	if err := r.Hooks.merge(other.Hooks); err != nil {
		return err
	}
	return nil
}

func createStacks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		stackLocation := location.appendPath(name)
		res[name] = Stack{
			Name:      name,
			Component: createComponentRef(env, stackLocation.appendPath("component"), yamlStack.Component, false),
			Hooks: StackHook{
				Deploy:   createHook(env, stackLocation.appendPath("hooks.deploy"), yamlStack.Hooks.Deploy),
				Undeploy: createHook(env, stackLocation.appendPath("hooks.undeploy"), yamlStack.Hooks.Undeploy)},
			  Parameters: createParameters(yamlStack.Params),
			  EnvVars:    createEnvVars(yamlStack.Env),
		}
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
