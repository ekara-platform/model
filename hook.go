package model

import (
	"encoding/json"
)

type Hook struct {
	Before []TaskRef
	After  []TaskRef
}

func (r Hook) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Before *[]TaskRef `json:",omitempty"`
		After  *[]TaskRef `json:",omitempty"`
	}{
		Before: &r.Before,
		After:  &r.After,
	})
}

func createHook(vErrs *ValidationErrors, location string, env *Environment, yamlHook yamlHook) Hook {
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}

	for i, yamlRef := range yamlHook.Before {
		hook.Before[i] = createTaskRef(vErrs, location+".before", env, yamlRef)
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i] = createTaskRef(vErrs, location+".after", env, yamlRef)
	}

	return hook
}
