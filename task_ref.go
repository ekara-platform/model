package model

import (
	"errors"
)

type (
	//TaskRef represents a reference to a task
	TaskRef struct {
		ref          string
		HookLocation hookLocation
		parameters   Parameters
		envVars      EnvVars
		env          *Environment
		location     DescriptorLocation
		mandatory    bool
	}
)

//reference return a validatable representation of the reference on a task
func (r TaskRef) validationDetails() refValidationDetails {
	result := make(map[string]interface{})
	for k, v := range r.env.Tasks {
		result[k] = v
	}
	return refValidationDetails{
		Id:        r.ref,
		Type:      "task",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

func (r *TaskRef) merge(other TaskRef) error {
	var err error
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters, err = r.parameters.inherit(other.parameters)
	if err != nil {
		return err
	}
	r.envVars, err = r.envVars.inherit(other.envVars)
	if err != nil {
		return err
	}
	return nil
}

// Resolve returns a resolved reference to a task containing all the
// inherited content from the referenced task
func (r TaskRef) Resolve() (Task, error) {
	var err error
	if err = ErrorOnInvalid(r); err.(ValidationErrors).HasErrors() {
		return Task{}, err
	}
	task := r.env.Tasks[r.ref]
	params, err := r.parameters.inherit(task.Parameters)
	if err != nil {
		return Task{}, err
	}
	envVars, err := r.envVars.inherit(task.EnvVars)
	if err != nil {
		return Task{}, err
	}
	return Task{
		Name:       task.Name,
		location:   task.location,
		cRef:       task.cRef,
		Playbook:   task.Playbook,
		Cron:       task.Cron,
		Hooks:      task.Hooks,
		Parameters: params,
		EnvVars:    envVars}, nil
}

func createTaskRef(env *Environment, location DescriptorLocation, tRef yamlTaskRef, hl hookLocation) (TaskRef, error) {
	params, err := CreateParameters(tRef.Params)
	if err != nil {
		return TaskRef{}, err
	}
	envVars, err := createEnvVars(tRef.Env)
	if err != nil {
		return TaskRef{}, err
	}
	return TaskRef{
		env:          env,
		HookLocation: hl,
		ref:          tRef.Task,
		parameters:   params,
		envVars:      envVars,
		location:     location,
		mandatory:    true,
	}, nil
}

func checkCircularRefs(taskRefs []TaskRef, alreadyEncountered *circularRefTracking) error {
	for _, taskRef := range taskRefs {
		if _, ok := (*alreadyEncountered)[taskRef.ref]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + taskRef.ref)
		}
	}
	return nil
}
