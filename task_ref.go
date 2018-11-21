package model

import (
	"errors"
)

type (
	TaskRef struct {
		ref        string
		parameters Parameters
		envVars    EnvVars
		env        *Environment
		location   DescriptorLocation
		mandatory  bool
	}
)

func (r TaskRef) Reference() Reference {
	result := make(map[string]interface{})
	for k, v := range r.env.Tasks {
		result[k] = v
	}
	return Reference{
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

func (r TaskRef) Resolve() (Task, error) {
	validationErrors := ErrorOn(r)
	if validationErrors.HasErrors() {
		return Task{}, validationErrors
	}
	task := r.env.Tasks[r.ref]
	return Task{
		Name:       task.Name,
		Parameters: r.parameters.inherits(task.Parameters),
		EnvVars:    r.envVars.inherits(task.EnvVars)}, nil
}

func createTaskRef(env *Environment, location DescriptorLocation, taskRef yamlTaskRef) TaskRef {
	if len(taskRef.Task) == 0 {
		env.errors.addError(errors.New("empty task reference"), location)
	} else {
		return TaskRef{
			env:        env,
			ref:        taskRef.Task,
			parameters: createParameters(taskRef.Params),
			envVars:    createEnvVars(taskRef.Env),
			location:   location,
			mandatory:  true,
		}
	}
	return TaskRef{}
}

func checkCircularRefs(taskRefs []yamlTaskRef, alreadyEncountered *circularRefTracking) error {
	for _, ref := range taskRefs {
		if _, ok := (*alreadyEncountered)[ref.Task]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + ref.Task)
		}
	}
	return nil
}
