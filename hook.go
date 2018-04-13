package model

type Hook struct {
	Before []TaskRef
	After  []TaskRef
}

func createHook(vErrs *ValidationErrors, tasks map[string]Task, location string, yamlHook yamlHook) Hook {
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}

	for i, yamlRef := range yamlHook.Before {
		hook.Before[i] = createTaskRef(vErrs, tasks, location+".before", yamlRef)
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i] = createTaskRef(vErrs, tasks, location+".after", yamlRef)
	}

	return hook
}
