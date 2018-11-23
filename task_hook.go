package model

import (
	"encoding/json"
)

type (
	//TaskHook represents hooks associated to a task
	TaskHook struct {
		//Execute specifies the hook tasks to run when a task is executed
		Execute Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r TaskHook) HasTasks() bool {
	return r.Execute.HasTasks()
}

func (r *TaskHook) merge(other TaskHook) error {
	if err := r.Execute.merge(other.Execute); err != nil {
		return err
	}
	return nil
}

func (r TaskHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Execute)
}

// MarshalJSON returns the serialized content of the hook as JSON
func (r TaskHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Execute *Hook `json:",omitempty"`
	}{}
	if r.Execute.HasTasks() {
		t.Execute = &r.Execute
	}
	return json.Marshal(t)
}
