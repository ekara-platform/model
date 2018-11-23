package model

import (
	"encoding/json"
)

type (
	//NodeHook represents hooks associated to a node set
	NodeHook struct {
		//Provisioned specifies the hook tasks to run when a node set is provisioned
		Provision Hook
		//Destroy specifies the hook tasks to run when a node set is destroyed
		Destroy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r NodeHook) HasTasks() bool {
	return r.Provision.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r *NodeHook) merge(other NodeHook) error {
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	if err := r.Destroy.merge(other.Destroy); err != nil {
		return err
	}
	return nil
}

func (r NodeHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Provision, r.Destroy)
}

// MarshalJSON returns the serialized content of the hook as JSON
func (r NodeHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Provision *Hook `json:",omitempty"`
		Destroy   *Hook `json:",omitempty"`
	}{}
	if r.Provision.HasTasks() {
		t.Provision = &r.Provision
	}
	if r.Destroy.HasTasks() {
		t.Destroy = &r.Destroy
	}
	return json.Marshal(t)
}
