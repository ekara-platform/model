package model

import (
	"errors"
	"fmt"
)

type (
	//Stack represent an Stack installable on the built environment
	Stack struct {
		// The component containing the stack
		cRef componentRef
		// The name of the stack
		Name string
		//DependsOn specifies the stack references on which this one depends
		DependsOn Dependencies
		// The hooks linked to the stack lifecycle events
		Hooks StackHook
		// The stack parameters
		Parameters Parameters
		// The stack environment variables
		EnvVars EnvVars
		// The stack content to be copied on volumes
		Copies Copies
	}

	//StackRef defines a dependency a on stack which must be previously processed
	StackRef struct {
		//ref defines the id of the referenced stack
		ref string
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
	// TODO add a validation on copies destinations
	if len(s.DependsOn.Content) > 0 {
		return ErrorOnInvalid(s.cRef, s.DependsOn.Content, s.Hooks)
	}
	return ErrorOnInvalid(s.cRef, s.Hooks)
}

func (s *Stack) merge(other Stack) error {
	var err error
	if s.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + s.Name + " != " + other.Name + ")")
	}
	if err = s.cRef.merge(other.cRef); err != nil {
		return err
	}
	s.Parameters, err = s.Parameters.inherit(other.Parameters)
	if err != nil {
		return err
	}
	s.EnvVars, err = s.EnvVars.inherit(other.EnvVars)
	if err != nil {
		return err
	}
	s.DependsOn = s.DependsOn.inherit(other.DependsOn)
	s.Copies = s.Copies.inherit(other.Copies)
	return s.Hooks.merge(other.Hooks)
}

func (r Stacks) merge(env *Environment, others Stacks) (Stacks, error) {
	res := make(map[string]Stack)
	for k, v := range r {
		res[k] = v
	}

	for id, s := range others {
		if stack, ok := res[id]; ok {
			sm := &stack
			if err := sm.merge(s); err != nil {
				return res, err
			}
			res[id] = *sm
		} else {
			s.cRef.env = env
			res[id] = s
		}
	}
	return res, nil
}

func createStacks(env *Environment, holder string, location DescriptorLocation, yamlEnv *yamlEnvironment) (Stacks, error) {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		// Root stack
		stackLocation := location.appendPath(name)
		params, err := CreateParameters(yamlStack.Params)
		if err != nil {
			return res, err
		}
		envVars, err := createEnvVars(yamlStack.Env)
		if err != nil {
			return res, err
		}
		dHook, err := createHook(env, stackLocation.appendPath("hooks.deploy"), yamlStack.Hooks.Deploy)
		if err != nil {
			return res, err
		}

		yC := yamlStack.Component
		if yC == "" || yC == "_" {
			yC = holder
		}

		s := Stack{
			Name: name,
			cRef: createComponentRef(env, stackLocation.appendPath("component"), yC, false),
			Hooks: StackHook{
				Deploy: dHook},
			Parameters: params,
			EnvVars:    envVars,
			DependsOn:  createDependencies(env, stackLocation.appendPath("depends_on"), name, yamlStack.DependsOn),
			Copies:     createCopies(env, stackLocation.appendPath("volume_copies"), yamlStack.Copies),
		}
		res[name] = s
		//env.Ekara.tagUsedComponent(res[name])
	}
	return res, nil
}

//Dependency returns the potential Stacks on which this one depends
func (s Stack) Dependency() (bool, []Stack) {
	res := make([]Stack, 0)
	for _, val := range s.DependsOn.Content {
		if val, ok := s.cRef.env.Stacks[val.ref]; ok {
			res = append(res, val)
		}
	}
	return len(res) > 0, res
}

//ResolveDependencies returns the stacks based on the order of the dependencies
func (r Stacks) ResolveDependencies() ([]Stack, error) {
	result := make([]Stack, 0, 0)
	if len(r) == 0 {
		return result, nil
	}

	g := newGraph(len(r))
	for _, vs := range r {
		g.addNode(vs.Name)
	}
	for _, vs := range r {
		for _, vd := range vs.DependsOn.Content {
			g.addEdge(vd.ref, vs.Name)
		}
	}
	res, ok := g.sort()
	if !ok {
		return result, fmt.Errorf("A cyclic dependency has been detected")
	}
	for _, val := range res {
		if stack, ok := r[val]; ok {
			result = append(result, stack)
		}
	}
	return result, nil
}

//validationDetails return a validatable representation of the reference on the stack
func (s StackRef) validationDetails() refValidationDetails {
	result := make(map[string]interface{})
	for k, v := range s.env.Stacks {
		result[k] = v
	}
	return refValidationDetails{
		Id:        s.ref,
		Type:      "stack dependency",
		Mandatory: true,
		Location:  s.location,
		Repo:      result,
	}
}

//Component returns the referenced component
func (s Stack) Component() (Component, error) {
	return s.cRef.resolve()
}

//ComponentName returns the referenced component name
func (s Stack) ComponentName() string {
	return s.cRef.ref
}

// Resolve returns a resolved reference to a stack
func (s StackRef) Resolve() (Stack, error) {
	var err error
	if err = ErrorOnInvalid(s); err.(ValidationErrors).HasErrors() {
		return Stack{}, err
	}
	task := s.env.Stacks[s.ref]
	return task, nil
}
