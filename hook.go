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
		//After specifies the tasks to run once the ekara life cycle event has occured
		After []TaskRef
	}
)

const (
	// Hook location before
	HOOK_BEFORE hookLocation = "Before"
	// Hook location after
	HOOK_AFTER hookLocation = "After"
)

func createHook(env *Environment, location DescriptorLocation, yamlHook yamlHook) Hook {
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}

	for i, yamlRef := range yamlHook.Before {
		hook.Before[i] = createTaskRef(env, location.appendPath("before"), yamlRef, HOOK_BEFORE)
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i] = createTaskRef(env, location.appendPath("after"), yamlRef, HOOK_AFTER)
	}
	return hook
}

func (r Hook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Before, r.After)
}

func (r *Hook) merge(other Hook) error {
	if !reflect.DeepEqual(r, &other) {
		r.Before = append(r.Before, other.Before...)
		r.After = append(r.After, other.After...)
	}
	return nil
}

//HasTasks returns true if the hook contains at least one task reference
func (r Hook) HasTasks() bool {
	return len(r.Before) > 0 || len(r.After) > 0
}
