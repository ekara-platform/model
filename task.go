package model

import (
	"bytes"
	"errors"
	"fmt"
)

type Task struct {
	root *Environment
	Labels
	Parameters

	Name     string
	Playbook string
	Cron     string
	RunOn    NodeSetRef

	Hooks struct {
		Execute Hook
	}
}

type TaskRef struct {
	Parameters
	task *Task
}

func createTasks(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Task {
	res := map[string]Task{}

	for name, yamlTask := range yamlEnv.Tasks {
		if len(yamlTask.Playbook) == 0 {
			vErrs.AddError(errors.New("missing playbook"), "tasks."+name+".playbook")
		}

		res[name] = Task{
			root:       env,
			Labels:     createLabels(vErrs, yamlTask.Labels...),
			Parameters: createParameters(vErrs, yamlTask.Params),
			Name:       name,
			Playbook:   yamlTask.Playbook,
			Cron:       yamlTask.Cron}
	}

	for name, yamlTask := range yamlEnv.Tasks {
		err := checkCircularRefs(yamlTask.Hooks.Execute.Before, &circularRefTracking{})
		if err != nil {
			vErrs.AddError(err, "tasks."+name+".hooks.execute.before")
		}
		err = checkCircularRefs(yamlTask.Hooks.Execute.After, &circularRefTracking{})
		if err != nil {
			vErrs.AddError(err, "tasks."+name+".hooks.execute.after")
		}
		task := res[name]
		task.Hooks.Execute = createHook(vErrs, res, "tasks."+name+".hooks.execute", yamlTask.Hooks.Execute)
		if len(yamlTask.RunOn) > 0 {
			task.RunOn = createNodeSetRef(vErrs, env, "tasks."+name+".runOn", yamlTask.RunOn...)
		}
	}

	return res
}

func createTaskRef(vErrs *ValidationErrors, tasks map[string]Task, location string, yamlRef yamlRef) TaskRef {
	if len(yamlRef.Name) == 0 {
		vErrs.AddError(errors.New("empty task reference"), location)
	} else {
		if val, ok := tasks[yamlRef.Name]; ok {
			return TaskRef{Parameters: createParameters(vErrs, yamlRef.Params), task: &val}
		} else {
			vErrs.AddError(errors.New("unknown task reference: "+yamlRef.Name), location)
		}
	}
	return TaskRef{}
}

func checkCircularRefs(yamlRefs []yamlRef, alreadyEncountered *circularRefTracking) error {
	for _, ref := range yamlRefs {
		if _, ok := (*alreadyEncountered)[ref.Name]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + ref.Name)
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
