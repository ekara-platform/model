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
		// The stack reference on which this one depends
		DependsOn StackDependency
		// The component containing the stack
		Component componentRef
		// The hooks linked to the stack lifecycle events
		Hooks      StackHook
		Parameters Parameters
		EnvVars    EnvVars
	}

	//StackDependency defines a dependency on a stack which must be previously processed
	StackDependency struct {
		Stack string
		//env specifies the environment holding the referenced component
		env *Environment
		//location indicates where the reference is located into the descriptor
		location DescriptorLocation
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
	if s.DependsOn.Stack != "" {
		return ErrorOnInvalid(s.Component, s.DependsOn, s.Hooks)
	}
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
func (s Stack) Dependency() (bool, Stack) {
	if val, ok := s.Component.env.Stacks[s.DependsOn.Stack]; ok {
		return true, val
	}
	return false, Stack{}
}

func createStacks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		stackLocation := location.appendPath(name)
		s := Stack{
			Name:      name,
			Component: createComponentRef(env, stackLocation.appendPath("component"), yamlStack.Component, false),
			Hooks: StackHook{
				Deploy:   createHook(env, stackLocation.appendPath("hooks.deploy"), yamlStack.Hooks.Deploy),
				Undeploy: createHook(env, stackLocation.appendPath("hooks.undeploy"), yamlStack.Hooks.Undeploy)},
			Parameters: createParameters(yamlStack.Params),
			EnvVars:    createEnvVars(yamlStack.Env),
		}
		if yamlStack.DependsOn != "" && yamlStack.DependsOn != name {
			depLocation := stackLocation.appendPath("depends_on")
			depLocation = depLocation.appendPath(yamlStack.DependsOn)
			s.DependsOn = StackDependency{
				Stack:    yamlStack.DependsOn,
				location: depLocation,
				env:      env,
			}
		}
		res[name] = s
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
				str := s.DependsOn.Stack
				if _, ok := done[str]; ok || str == "" {
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

//reference return a validatable representation of the reference on the component
func (s StackDependency) reference() validatableReference {
	result := make(map[string]interface{})
	for k, v := range s.env.Stacks {
		result[k] = v
	}
	return validatableReference{
		Id:        s.Stack,
		Type:      "stack dependency",
		Mandatory: true,
		Location:  s.location,
		Repo:      result,
	}
}
