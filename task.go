package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	Task struct {
		// Name of the task
		Name string
		// The component containing the task
		Component ComponentRef
		// The playbook to execute
		Playbook string
		// The cron expression when the task must be scheduled
		Cron string
		// The task parameters
		Parameters Parameters
		// The task environment variables
		EnvVars EnvVars
		// Hooks for executing other tasks around execution
		Hooks TaskHook
	}

	Tasks map[string]Task

	TaskHook struct {
		Execute Hook
	}

	TaskRef struct {
		ref        string
		parameters Parameters
		envVars    EnvVars

		env      *Environment
		location DescriptorLocation
	}

	circularRefTracking map[string]interface{}
)

func (r TaskHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Execute *Hook `json:",omitempty"`
	}{}
	if r.Execute.HasTasks() {
		t.Execute = &r.Execute
	}
	return json.Marshal(t)
}

func (r *TaskHook) merge(other TaskHook) error {
	if err := r.Execute.merge(other.Execute); err != nil {
		return err
	}
	return nil
}

func (r TaskHook) HasTasks() bool {
	return r.Execute.HasTasks()
}

func (r Task) MarshalJSON() ([]byte, error) {
	t := struct {
		Name       string      `json:",omitempty"`
		Playbook   string      `json:",omitempty"`
		Cron       string      `json:",omitempty"`
		On         []string    `json:",omitempty"`
		Parameters *Parameters `json:",omitempty"`
		EnvVars    *EnvVars    `json:",omitempty"`
		Hooks      *TaskHook   `json:",omitempty"`
	}{
		Name:       r.Name,
		Playbook:   r.Playbook,
		Cron:       r.Cron,
		Parameters: &r.Parameters,
		EnvVars:    &r.EnvVars,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

func (r Task) validate() ValidationErrors {
	vErrs := r.Component.validate()
	vErrs.merge(r.Hooks.Execute.validate())
	return vErrs
}

func (r *Task) merge(other Task) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Component.merge(other.Component); err != nil {
		return err
	}
	if err := r.Hooks.merge(other.Hooks); err != nil {
		return err
	}
	if r.Playbook == "" {
		r.Playbook = other.Playbook
	}
	if r.Cron == "" {
		r.Cron = other.Cron
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
	return nil
}

func createTasks(env *Environment, yamlEnv *yamlEnvironment) Tasks {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		if len(yamlTask.Playbook) == 0 {
			env.errors.addError(errors.New("empty playbook path"), env.location.appendPath("tasks."+name+".playbook"))
		}
		err := checkCircularRefs(yamlTask.Hooks.Execute.Before, &circularRefTracking{})
		if err != nil {
			env.errors.addError(err, env.location.appendPath("tasks."+name+".hooks.execute.before"))
		}
		err = checkCircularRefs(yamlTask.Hooks.Execute.After, &circularRefTracking{})
		if err != nil {
			env.errors.addError(err, env.location.appendPath("tasks."+name+".hooks.execute.after"))
		}

		res[name] = Task{
			Name:       name,
			Playbook:   yamlTask.Playbook,
			Component:  createComponentRef(env, env.location.appendPath("tasks."+name+".component"), yamlTask.Component, false),
			Cron:       yamlTask.Cron,
			Parameters: createParameters(yamlTask.Params),
			EnvVars:    createEnvVars(yamlTask.Env),
			Hooks: TaskHook{
				Execute: createHook(env, env.location.appendPath("tasks."+name+".hooks.execute"), yamlTask.Hooks.Execute),
			},
		}
	}
	return res
}

func (r Tasks) merge(env *Environment, other Tasks) error {
	for id, t := range other {
		if task, ok := r[id]; ok {
			if err := task.merge(t); err != nil {
				return err
			}
		} else {
			t.Component.env = env
			r[id] = t
		}
	}
	return nil
}

func (r TaskRef) Resolve() (Task, error) {
	validationErrors := r.validate()
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
			location:   location}
	}
	return TaskRef{}
}

func (r TaskRef) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.ref) == 0 {
		vErrs.addError(errors.New("empty task reference"), r.location)
	} else {
		if _, ok := r.env.Providers[r.ref]; !ok {
			vErrs.addError(errors.New("reference to unknown task: "+r.ref), r.location)
		}
	}
	return vErrs
}

func (r *TaskRef) merge(other TaskRef) {
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
}

func checkCircularRefs(taskRefs []yamlTaskRef, alreadyEncountered *circularRefTracking) error {
	for _, ref := range taskRefs {
		if _, ok := (*alreadyEncountered)[ref.Task]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + ref.Task)
		}
	}
	return nil
}

func (r *circularRefTracking) String() string {
	b := new(bytes.Buffer)
	for key := range *r {
		fmt.Fprintf(b, "%s -> ", key)
	}
	return b.String()
}
