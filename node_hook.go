package model

import (
	"encoding/json"
)

type (
	NodeHook struct {
		Provision Hook
		Destroy   Hook
	}
)

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
	vErrs := ValidationErrors{}
	vErrs.merge(ErrorOn(r.Provision))
	vErrs.merge(ErrorOn(r.Destroy))
	return vErrs
}

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
