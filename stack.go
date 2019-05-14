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
		DependsOn []stackRef
		// The hooks linked to the stack lifecycle events
		Hooks      StackHook
		Parameters Parameters
		EnvVars    EnvVars
	}

	//stackRef defines a dependency a on stack which must be previously processed
	stackRef struct {
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
	if len(s.DependsOn) > 0 {
		return ErrorOnInvalid(s.cRef, s.DependsOn, s.Hooks)
	}
	return ErrorOnInvalid(s.cRef, s.Hooks)
}

func (s *Stack) merge(other Stack) error {
	if s.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + s.Name + " != " + other.Name + ")")
	}
	if err := s.cRef.merge(other.cRef); err != nil {
		return err
	}
	s.Parameters = s.Parameters.inherits(other.Parameters)
	s.EnvVars = s.EnvVars.inherits(other.EnvVars)

	return s.Hooks.merge(other.Hooks)
}

//Dependency returns the potential Stacks which this one depends
func (s Stack) Dependency() (bool, []Stack) {
	res := make([]Stack, 0)
	for _, val := range s.DependsOn {
		if val, ok := s.cRef.env.Stacks[val.ref]; ok {
			res = append(res, val)
		}
	}
	return len(res) > 0, res
}

func createStacks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		stackLocation := location.appendPath(name)
		s := Stack{
			Name: name,
			cRef: createComponentRef(env, stackLocation.appendPath("component"), yamlStack.Component, false),
			Hooks: StackHook{
				Deploy:   createHook(env, stackLocation.appendPath("hooks.deploy"), yamlStack.Hooks.Deploy),
				Undeploy: createHook(env, stackLocation.appendPath("hooks.undeploy"), yamlStack.Hooks.Undeploy)},
			Parameters: createParameters(yamlStack.Params),
			EnvVars:    createEnvVars(yamlStack.Env),
		}
		deps := make([]stackRef, 0, 0)
		for _, v := range yamlStack.DependsOn {
			if v != "" && v != name {
				depLocation := stackLocation.appendPath("depends_on")
				depLocation = depLocation.appendPath(v)
				dep := stackRef{
					ref:      v,
					location: depLocation,
					env:      env,
				}
				deps = append(deps, dep)
			}
		}
		if len(deps) > 0 {
			s.DependsOn = deps
		}
		res[name] = s
		env.Ekara.tagUsedComponent(res[name])
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
			s.cRef.env = env
			r[id] = s
		}
	}
	return nil
}

//ResolveDependencies returns a channel to get access
//to the stacks based on the rorder of the dependencies
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
		for _, vd := range vs.DependsOn {
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
func (s stackRef) validationDetails() refValidationDetails {
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

func (r Stack) Component() (Component, error) {
	return r.cRef.resolve()
}

func (r Stack) ComponentName() string {
	return r.cRef.ref
}
