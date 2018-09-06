package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Task struct {
	// Name of the task
	Name string
	// The playbook to execute
	Playbook string
	// The cron expression when the task must be scheduled
	Cron string
	// The nodes to run the task on
	On NodeSetRef
	// The task parameters
	Parameters Parameters
	// The task environment variables
	EnvVars EnvVars
	// Hooks for executing other tasks around execution
	Hooks TaskHook
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

	for _, k := range r.On.nodeSets {
		t.On = append(t.On, k.Name)
	}

	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

type TaskHook struct {
	Execute Hook
}

func (r TaskHook) HasTasks() bool {
	return r.Execute.HasTasks()
}

func (r TaskHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Execute *Hook `json:",omitempty"`
	}{}
	if r.Execute.HasTasks() {
		t.Execute = &r.Execute
	}
	return json.Marshal(t)
}

type TaskRef struct {
	task       *Task
	parameters Parameters
	envVars    EnvVars
}

func createTasks(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Task {
	res := map[string]Task{}
	for name, yamlTask := range yamlEnv.Tasks {
		if len(yamlTask.Playbook) == 0 {
			vErrs.AddError(errors.New("missing playbook"), "tasks."+name+".playbook")
		}
		err := checkCircularRefs(yamlTask.Hooks.Execute.Before, &circularRefTracking{})
		if err != nil {
			vErrs.AddError(err, "tasks."+name+".hooks.execute.before")
		}
		err = checkCircularRefs(yamlTask.Hooks.Execute.After, &circularRefTracking{})
		if err != nil {
			vErrs.AddError(err, "tasks."+name+".hooks.execute.after")
		}

		res[name] = Task{
			Name:       name,
			Playbook:   yamlTask.Playbook,
			Cron:       yamlTask.Cron,
			On:         createNodeSetRef(vErrs, env, "tasks."+name+".on", yamlTask.On...),
			Parameters: createParameters(yamlTask.Params),
			EnvVars:    createEnvVars(yamlTask.Env)}
	}
	return res
}

// TODO Add units tests for this on the "complete_descriptor"
func createTaskRef(vErrs *ValidationErrors, location string, env *Environment, taskRef yamlTaskRef) TaskRef {
	if len(taskRef.Task) == 0 {
		vErrs.AddError(errors.New("empty task reference"), location)
	} else {
		if val, ok := env.Tasks[taskRef.Task]; ok {
			return TaskRef{
				task:       &val,
				parameters: createParameters(taskRef.Params).inherit(val.Parameters),
				envVars:    createEnvVars(taskRef.Env).inherit(val.EnvVars),
			}
		} else {
			vErrs.AddError(errors.New("unknown task reference: "+taskRef.Task), location)
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

type circularRefTracking map[string]interface{}

func (r *circularRefTracking) String() string {
	b := new(bytes.Buffer)
	for key := range *r {
		fmt.Fprintf(b, "%s -> ", key)
	}
	return b.String()
}
