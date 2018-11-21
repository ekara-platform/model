package model

type Hook struct {
	Before []TaskRef
	After  []TaskRef
}

func createHook(env *Environment, location DescriptorLocation, yamlHook yamlHook) Hook {
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}

	for i, yamlRef := range yamlHook.Before {
		hook.Before[i] = createTaskRef(env, location.appendPath("before"), yamlRef)
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i] = createTaskRef(env, location.appendPath("after"), yamlRef)
	}

	return hook
}

func (r Hook) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(ErrorOn(r.Before))
	vErrs.merge(ErrorOn(r.After))
	return vErrs
}

func (r *Hook) merge(other Hook) error {
	r.Before = append(r.Before, other.Before...)
	r.After = append(other.After, r.After...)
	return nil
}

func (r Hook) HasTasks() bool {
	return len(r.Before) > 0 || len(r.After) > 0
}
