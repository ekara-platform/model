package descriptor

type Hook struct {
	Before []TaskRef
	After  []TaskRef
}

func createHook(tasks map[string]Task, yamlHook yamlHook) (res Hook, err error) {
	hook := Hook{
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}

	for i, yamlRef := range yamlHook.Before {
		hook.Before[i], err = createTaskRef(tasks, yamlRef)
		if err != nil {
			return
		}
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i], err = createTaskRef(tasks, yamlRef)
		if err != nil {
			return
		}
	}

	return
}
