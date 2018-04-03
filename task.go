package descriptor

import (
	"errors"
	"bytes"
	"fmt"
)

type Task struct {
	root *Environment
	Labels
	Parameters

	Name     string
	Playbook string
	Cron     string

	Hooks struct {
		Execute Hook
	}
}

type TaskRef struct {
	Parameters
	task *Task
}

func createTasks(env *Environment, yamlEnv *yamlEnvironment) (res map[string]Task, err error) {
	res = map[string]Task{}
	for name, yamlTask := range yamlEnv.Tasks {
		res[name] = Task{
			root:       env,
			Labels:     createLabels(yamlTask.Labels...),
			Parameters: createParameters(yamlTask.Params),
			Name:       name,
			Playbook:   yamlTask.Playbook,
			Cron:       yamlTask.Cron}
	}

	for name, yamlTask := range yamlEnv.Tasks {
		err = checkCircularRefs(yamlTask.Hooks.Execute.Before, &circularRefTracking{})
		if err != nil {
			return
		}
		err = checkCircularRefs(yamlTask.Hooks.Execute.After, &circularRefTracking{})
		if err != nil {
			return
		}
		task := res[name]
		task.Hooks.Execute, err = createHook(res, yamlTask.Hooks.Execute)
	}

	return
}

func createTaskRef(tasks map[string]Task, yamlRef yamlRef) (res TaskRef, err error) {
	if val, ok := tasks[yamlRef.Name]; ok {
		res = TaskRef{Parameters: createParameters(yamlRef.Params), task: &val}
	} else {
		err = errors.New("unknown task reference: " + yamlRef.Name)
	}
	return
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
