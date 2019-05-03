package model

import (
	"errors"
	"sort"
)

type (
	//Stack represent an Stack installable on the built environment
	Stack struct {
		// The name of the stack
		Name string
		// The name of the stack on which this one depends
		DependsOn string
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
func (s Stack) DescType() string {
	return "Stack"
}

//DescName returns the Describable name of the stack
func (s Stack) DescName() string {
	return s.Name
}

func (s Stack) validate() ValidationErrors {
	return ErrorOnInvalid(s.Component, s.Hooks)
}

func (s *Stack) merge(other Stack) error {
	if s.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + s.Name + " != " + other.Name + ")")
	}
	if err := s.Component.merge(other.Component); err != nil {
		return err
	}
	s.Parameters = s.Parameters.inherits(other.Parameters)
	s.EnvVars = s.EnvVars.inherits(other.EnvVars)

	return s.Hooks.merge(other.Hooks)
}

//Dependency returns the potential Stack which this one depends
func (r Stack) Dependency() (bool, Stack) {
	if val, ok := r.Component.env.Stacks[r.DependsOn]; ok {
		return true, val
	}
	return false, Stack{}
}

func createStacks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		stackLocation := location.appendPath(name)
		res[name] = Stack{
			Name:      name,
			DependsOn: yamlStack.DependsOn,
			Component: createComponentRef(env, stackLocation.appendPath("component"), yamlStack.Component, false),
			Hooks: StackHook{
				Deploy:   createHook(env, stackLocation.appendPath("hooks.deploy"), yamlStack.Hooks.Deploy),
				Undeploy: createHook(env, stackLocation.appendPath("hooks.undeploy"), yamlStack.Hooks.Undeploy)},
			Parameters: createParameters(yamlStack.Params),
			EnvVars:    createEnvVars(yamlStack.Env),
		}
		env.Ekara.tagUsedComponent(res[name].Component)
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

//ResolveDependencies returns a channel to get access
//to the stacks based on the rorder of the dependencies
func (r Stacks) ResolveDependencies() <-chan Stack {
	out := make(chan Stack)
	go func() {
		// The stacks to process
		todo := Stacks{}
		// The stacks already processed
		done := Stacks{}
		//We will work on a copy in order to leave the original Stacks untouched
		for k, val := range r {
			todo[k] = val
		}
		//We still have stuff to process
		for len(todo) > 0 {
			var keys []string
			for k := range todo {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				s := todo[k]
				//We can return the stack if it depends on nothing
				//or if the dependency has already been processed
				if _, ok := done[s.DependsOn]; ok || s.DependsOn == "" {
					out <- s
					done[k] = s
					delete(todo, k)
					continue
				}
			}
		}
		//All done, bye...
		close(out)
	}()
	return out
}
