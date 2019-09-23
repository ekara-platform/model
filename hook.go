package model

import (
	"reflect"
)

type (
	// Type used to identify a hook location, typically 'before' or 'after'
	hookLocation string

	// Hook represents tasks to be executed linked to an ekara life cycle event
	Hook struct {
		//Before specifies the tasks to run before the ekara life cycle event occurs
		Before []TaskRef
		//After specifies the tasks to run once the ekara life cycle event has occurred
		After []TaskRef
	}
)

const (
	//HookBefore Hook located before a task
	HookBefore hookLocation = "Before"
	//HookAfter Hook located after a task
	HookAfter hookLocation = "After"
)

func createHook(env *Environment, location DescriptorLocation, yamlHook yamlHook) (Hook, error) {
	var err error
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}
	for i, yamlRef := range yamlHook.Before {
		hook.Before[i], err = createTaskRef(env, location.appendPath("before"), yamlRef, HookBefore)
		if err != nil {
			return hook, err
		}
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i], err = createTaskRef(env, location.appendPath("after"), yamlRef, HookAfter)
		if err != nil {
			return hook, err
		}
	}
	return hook, nil
}

func (r Hook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Before, r.After)
}

func (r *Hook) customize(with Hook) error {
	if !reflect.DeepEqual(r, &with) {
		r.Before = append(r.Before, with.Before...)
		r.After = append(r.After, with.After...)
	}
	return nil
}

//HasTasks returns true if the hook contains at least one task reference
func (r Hook) HasTasks() bool {
	return len(r.Before) > 0 || len(r.After) > 0
}
