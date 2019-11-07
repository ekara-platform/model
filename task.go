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
		Playbook string `yaml:",omitempty"`
		// The cron expression when the task must be scheduled
		Cron string `yaml:",omitempty"`
		// The task parameters
		Parameters Parameters `yaml:",omitempty"`
		// The task environment variables
		EnvVars EnvVars `yaml:",omitempty"`
		//The hooks linked to the task lifecycle events
		Hooks TaskHook `yaml:",omitempty"`
	}

	//Tasks represent all the tasks of an environment
	Tasks map[string]*Task

	circularRefTracking map[string]interface{}
)

//DescType returns the Describable type of the task
//  Hardcoded to : "Task"
func (r Task) DescType() string {
	return "Task"
}

//DescName returns the Describable name of the task
func (r Task) DescName() string {
	return r.Name
}

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

func (r *Task) customize(with Task) error {
	var err error
	if !reflect.DeepEqual(r, &with) {
		if r.Name != with.Name {
			return errors.New("cannot customize unrelated tasks (" + r.Name + " != " + with.Name + ")")
		}
		if err = r.cRef.customize(with.cRef); err != nil {
			return err
		}
		if err = r.Hooks.customize(with.Hooks); err != nil {
			return err
		}

		r.Playbook = with.Playbook
		r.Cron = with.Cron

		r.Parameters = with.Parameters.inherit(r.Parameters)
		r.EnvVars = with.EnvVars.inherit(r.EnvVars)
	}
	return nil
}

func createTasks(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Tasks, error) {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		taskLocation := location.appendPath(name)
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
			Parameters: CreateParameters(yamlTask.Params),
			EnvVars:    createEnvVars(yamlTask.Env),
			Hooks: TaskHook{
				Execute: eHook,
			},
		}
		//env.Ekara.tagUsedComponent(res[name])
	}
	return res, nil
}

func (r Tasks) customize(env *Environment, with Tasks) (Tasks, error) {

	work := make(map[string]*Task)
	for k, v := range r {
		work[k] = v
	}

	for id, t := range with {
		if task, ok := work[id]; ok {
			if err := task.customize(*t); err != nil {
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
