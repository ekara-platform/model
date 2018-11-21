package model

import (
	"encoding/json"
)

type (
	EnvironmentHooks struct {
		Init      Hook
		Provision Hook
		Deploy    Hook
		Undeploy  Hook
		Destroy   Hook
	}
)

func (r EnvironmentHooks) HasTasks() bool {
	return r.Init.HasTasks() ||
		r.Provision.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r *EnvironmentHooks) merge(other EnvironmentHooks) error {
	if err := r.Init.merge(other.Init); err != nil {
		return err
	}
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	if err := r.Deploy.merge(other.Deploy); err != nil {
		return err
	}
	if err := r.Undeploy.merge(other.Undeploy); err != nil {
		return err
	}
	if err := r.Destroy.merge(other.Destroy); err != nil {
		return err
	}
	return nil
}

func (r EnvironmentHooks) validate() ValidationErrors {
	return ErrorOnInvalid(r.Init, r.Provision, r.Deploy, r.Undeploy, r.Destroy)
}

func (r EnvironmentHooks) MarshalJSON() ([]byte, error) {
	t := struct {
		Init      *Hook `json:",omitempty"`
		Provision *Hook `json:",omitempty"`
		Deploy    *Hook `json:",omitempty"`
		Undeploy  *Hook `json:",omitempty"`
		Destroy   *Hook `json:",omitempty"`
	}{}

	if r.Init.HasTasks() {
		t.Init = &r.Init
	}
	if r.Provision.HasTasks() {
		t.Provision = &r.Provision
	}
	if r.Deploy.HasTasks() {
		t.Deploy = &r.Deploy
	}
	if r.Undeploy.HasTasks() {
		t.Undeploy = &r.Undeploy
	}
	if r.Destroy.HasTasks() {
		t.Destroy = &r.Destroy
	}

	return json.Marshal(t)
}
