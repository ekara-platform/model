package model

import (
	"errors"
)

type (
	taskRef struct {
		ref        string
		parameters Parameters
		envVars    EnvVars
		env        *Environment
		location   DescriptorLocation
		mandatory  bool
	}
)

//reference return a validatable representation of the reference on a task
func (r taskRef) reference() validatableReference {
	result := make(map[string]interface{})
	for k, v := range r.env.Tasks {
		result[k] = v
	}
	return validatableReference{
		Id:        r.ref,
		Type:      "task",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

func (r *taskRef) merge(other taskRef) {
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
}

func (r taskRef) Resolve() (Task, error) {
	validationErrors := ErrorOnInvalid(r)
	if validationErrors.HasErrors() {
		return Task{}, validationErrors
	}
	task := r.env.Tasks[r.ref]
	return Task{
		Name:       task.Name,
		Parameters: r.parameters.inherits(task.Parameters),
		EnvVars:    r.envVars.inherits(task.EnvVars)}, nil
}

func createTaskRef(env *Environment, location DescriptorLocation, tRef yamlTaskRef) taskRef {
	if len(tRef.Task) == 0 {
		env.errors.addError(errors.New("empty task reference"), location)
	} else {
		return taskRef{
			env:        env,
			ref:        tRef.Task,
			parameters: createParameters(tRef.Params),
			envVars:    createEnvVars(tRef.Env),
			location:   location,
			mandatory:  true,
		}
	}
	return taskRef{}
}

func checkCircularRefs(taskRefs []yamlTaskRef, alreadyEncountered *circularRefTracking) error {
	for _, ref := range taskRefs {
		if _, ok := (*alreadyEncountered)[ref.Task]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + ref.Task)
		}
	}
	return nil
}
