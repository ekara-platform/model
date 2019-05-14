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

func (r *TaskRef) merge(other TaskRef) {
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
}

// Resolve returns a resolved reference to a task containing all the
// inherited content from the referenced task
func (r TaskRef) Resolve() (Task, error) {
	validationErrors := ErrorOnInvalid(r)
	if validationErrors.HasErrors() {
		return Task{}, validationErrors
	}

	task := r.env.Tasks[r.ref]
	return Task{
		Name:       task.Name,
		location:   task.location,
		cRef:       task.cRef,
		Playbook:   task.Playbook,
		Cron:       task.Cron,
		Hooks:      task.Hooks,
		Parameters: r.parameters.inherits(task.Parameters),
		EnvVars:    r.envVars.inherits(task.EnvVars)}, nil
}

func createTaskRef(env *Environment, location DescriptorLocation, tRef yamlTaskRef, hl hookLocation) TaskRef {
	return TaskRef{
		env:          env,
		HookLocation: hl,
		ref:          tRef.Task,
		parameters:   createParameters(tRef.Params),
		envVars:      createEnvVars(tRef.Env),
		location:     location,
		mandatory:    true,
	}
}

func checkCircularRefs(taskRefs []TaskRef, alreadyEncountered *circularRefTracking) error {
	for _, taskRef := range taskRefs {
		if _, ok := (*alreadyEncountered)[taskRef.ref]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + taskRef.ref)
		}
	}
	return nil
}
