package model

import (
	"encoding/json"
)

type (
	TaskHook struct {
		Execute Hook
	}
)

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

func (r TaskHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Execute *Hook `json:",omitempty"`
	}{}
	if r.Execute.HasTasks() {
		t.Execute = &r.Execute
	}
	return json.Marshal(t)
}
