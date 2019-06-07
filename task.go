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
		// The component containing the task
		cRef     componentRef
		location DescriptorLocation
		// Name of the task
		Name string
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
	var err error
	if !reflect.DeepEqual(r, &other) {
		if r.Name != other.Name {
			return errors.New("cannot merge unrelated tasks (" + r.Name + " != " + other.Name + ")")
		}
		if err = r.cRef.merge(other.cRef); err != nil {
			return err
		}
		if err = r.Hooks.merge(other.Hooks); err != nil {
			return err
		}
		if r.Playbook == "" {
			r.Playbook = other.Playbook
		}
		if r.Cron == "" {
			r.Cron = other.Cron
		}
		r.Parameters, err = r.Parameters.inherit(other.Parameters)
		if err != nil {
			return err
		}
		r.EnvVars, err = r.EnvVars.inherit(other.EnvVars)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTasks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Tasks, error) {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		taskLocation := location.appendPath(name)
		params, err := CreateParameters(yamlTask.Params)
		if err != nil {
			return res, err
		}
		envVars, err := createEnvVars(yamlTask.Env)
		if err != nil {
			return res, err
		}
		eHook, err := createHook(env, taskLocation.appendPath("hooks.execute"), yamlTask.Hooks.Execute)
		if err != nil {
			return res, err
		}
		res[name] = &Task{
			location:   taskLocation,
			Name:       name,
			Playbook:   yamlTask.Playbook,
			cRef:       createComponentRef(env, taskLocation.appendPath("component"), yamlTask.Component, false),
			Cron:       yamlTask.Cron,
			Parameters: params,
			EnvVars:    envVars,
			Hooks: TaskHook{
				Execute: eHook,
			},
		}
		env.Ekara.tagUsedComponent(res[name])
	}
	return res, nil
}

func (r Tasks) merge(env *Environment, other Tasks) (Tasks, error) {

	work := make(map[string]*Task)
	for k, v := range r {
		work[k] = v
	}

	for id, t := range other {
		if task, ok := work[id]; ok {
			if err := task.merge(*t); err != nil {
				return work, err
			}
		} else {
			t.cRef.env = env
			work[id] = t
		}
	}
	return work, nil
}

func (r *circularRefTracking) String() string {
	b := new(bytes.Buffer)
	for key := range *r {
		fmt.Fprintf(b, "%s -> ", key)
	}
	return b.String()
}

//Component returns the referenced component
func (r Task) Component() (Component, error) {
	return r.cRef.resolve()
}

//ComponentName returns the referenced component name
func (r Task) ComponentName() string {
	return r.cRef.ref
}
