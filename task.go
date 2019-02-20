package model

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

type (
	//Task represent an task executable on the built environment
	Task struct {
		location DescriptorLocation
		// Name of the task
		Name string
		// The component containing the task
		Component componentRef
		// The playbook to execute
		Playbook string
		// The cron expression when the task must be scheduled
		Cron string
		// The task parameters
		Parameters Parameters
		// The task environment variables
		EnvVars EnvVars
		//The hooks linked to the task lifecycle events
		Hooks TaskHook
	}

	//Tasks represent all the tasks of an environment
	Tasks map[string]*Task

	circularRefTracking map[string]interface{}
)

func (r Task) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.Playbook) == 0 {
		vErrs.addError(errors.New("empty playbook path"), r.location.appendPath("playbook"))
	}
	err := checkCircularRefs(r.Hooks.Execute.Before, &circularRefTracking{})
	if err != nil {
		vErrs.addError(err, r.location.appendPath("hooks.execute.before"))
	}
	err = checkCircularRefs(r.Hooks.Execute.After, &circularRefTracking{})
	if err != nil {
		vErrs.addError(err, r.location.appendPath("hooks.execute.after"))
	}
	vErrs.merge(ErrorOnInvalid(r.Component, r.Hooks))
	return vErrs
}

func (r *Task) merge(other Task) error {
	if !reflect.DeepEqual(r, &other) {
		if r.Name != other.Name {
			return errors.New("cannot merge unrelated tasks (" + r.Name + " != " + other.Name + ")")
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
	}
	return nil
}

func createTasks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Tasks {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		taskLocation := location.appendPath(name)
		res[name] = &Task{
			location:   taskLocation,
			Name:       name,
			Playbook:   yamlTask.Playbook,
			Component:  createComponentRef(env, taskLocation.appendPath("component"), yamlTask.Component, false),
			Cron:       yamlTask.Cron,
			Parameters: createParameters(yamlTask.Params),
			EnvVars:    createEnvVars(yamlTask.Env),
			Hooks: TaskHook{
				Execute: createHook(env, taskLocation.appendPath("hooks.execute"), yamlTask.Hooks.Execute),
			},
		}
	}
	return res
}

func (r Tasks) merge(env *Environment, other Tasks) error {
	for id, t := range other {
		if task, ok := r[id]; ok {
			if err := task.merge(*t); err != nil {
				return err
			}
		} else {
			t.Component.env = env
			r[id] = t
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
