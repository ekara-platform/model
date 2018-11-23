package model

import (
	"encoding/json"
)

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Init specifies the hook tasks to run at the environment initialization
		Init Hook
		//Provisione specifies the hook tasks to run when the environment is provisioned
		Provision Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
		//Undeploy specifies the hook tasks to run when the environment is undeployed
		Undeploy Hook
		//Destroy specifies the hook tasks to run when the environment is destroyed
		Destroy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
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

// MarshalJSON returns the serialized content of the hook as JSON
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
