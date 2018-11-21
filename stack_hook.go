package model

import (
	"encoding/json"
)

type (
	StackHook struct {
		Deploy   Hook
		Undeploy Hook
	}
)

func (r StackHook) HasTasks() bool {
	return r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks()
}

func (r StackHook) validate() ValidationErrors {
	return ErrorOn(r.Deploy, r.Undeploy)
}

func (r StackHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Deploy   *Hook `json:",omitempty"`
		Undeploy *Hook `json:",omitempty"`
	}{}

	if r.Deploy.HasTasks() {
		t.Deploy = &r.Deploy
	}
	if r.Undeploy.HasTasks() {
		t.Undeploy = &r.Undeploy
	}
	return json.Marshal(t)
}
