package model

import (
	"encoding/json"
)

type (
	//StackHook represents hooks associated to a task
	StackHook struct {
		//Deploy specifies the hook tasks to run when a stack is deployed
		Deploy Hook
		//Undeploy specifies the hook tasks to run when a stack is undeployed
		Undeploy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r StackHook) HasTasks() bool {
	return r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks()
}

func (r StackHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Deploy, r.Undeploy)
}

// MarshalJSON returns the serialized content of the hook as JSON
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
