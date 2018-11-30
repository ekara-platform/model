package model

import (
	"encoding/json"
)

type (
	//Orchestrator specifies the orchestrator used to manage the environemt
	Orchestrator struct {
		// The component containing the orchestrator
		Component componentRef
		// The orchestrator parameters
		Parameters Parameters
		// The Docker parameters
		Docker Parameters
		// The orchestrator environment variables
		EnvVars EnvVars
	}
)

func createOrchestrator(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	return Orchestrator{
		Component:  createComponentRef(env, location.appendPath("component"), yamlO.Component, true),
		Parameters: createParameters(yamlO.Params),
		Docker:     createParameters(yamlO.Docker),
		EnvVars:    createEnvVars(yamlO.Env)}
}

func (r Orchestrator) validate() ValidationErrors {
	return ErrorOnInvalid(r.Component)
}

func (r *Orchestrator) merge(other Orchestrator) error {
	if err := r.Component.merge(other.Component); err != nil {
		return err
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.Docker = r.Docker.inherits(other.Docker)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
	return nil
}

// MarshalJSON returns the serialized content of orchestator as JSON
func (r Orchestrator) MarshalJSON() ([]byte, error) {
	component, e := r.Component.Resolve()
	if e != nil {
		return nil, e
	}
	return json.Marshal(struct {
		Component  string     `json:",omitempty"`
		Parameters Parameters `json:",omitempty"`
		Docker     Parameters `json:",omitempty"`
		EnvVars    EnvVars    `json:",omitempty"`
	}{
		Component:  component.Id,
		Parameters: r.Parameters,
		Docker:     r.Docker,
		EnvVars:    r.EnvVars,
	})
}
