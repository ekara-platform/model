package model

import (
	"encoding/json"
)

type (
	Orchestrator struct {
		// The component containing the orchestrator
		Component ComponentRef
		// The orchestrator parameters
		Parameters Parameters
		// The Docker parameters
		Docker Parameters
		// The orchestrator environment variables
		EnvVars EnvVars
	}

	OrchestratorRef struct {
		parameters Parameters
		docker     Parameters
		envVars    EnvVars

		env      *Environment
		location DescriptorLocation
	}
)

func createOrchestrator(env *Environment, yamlEnv *yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	return Orchestrator{
		Component:  createComponentRef(env, env.location.appendPath("orchestrator"), yamlO.Component, true),
		Parameters: createParameters(yamlO.Params),
		Docker:     createParameters(yamlO.Docker),
		EnvVars:    createEnvVars(yamlO.Env)}
}

func (r Orchestrator) validate() ValidationErrors {
	return r.Component.validate()
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

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) OrchestratorRef {
	return OrchestratorRef{
		env:        env,
		parameters: createParameters(yamlRef.Params),
		docker:     createParameters(yamlRef.Docker),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location}
}

func (r OrchestratorRef) validate() ValidationErrors {
	return ValidationErrors{}
}

func (r *OrchestratorRef) merge(other OrchestratorRef) error {
	return nil
}

func (r OrchestratorRef) Resolve() (Orchestrator, error) {
	validationErrors := r.validate()
	if validationErrors.HasErrors() {
		return Orchestrator{}, validationErrors
	}
	orchestrator := r.env.Orchestrator
	return Orchestrator{
		Component:  orchestrator.Component,
		Parameters: r.parameters.inherits(orchestrator.Parameters),
		Docker:     r.docker.inherits(orchestrator.Docker),
		EnvVars:    r.envVars.inherits(orchestrator.EnvVars)}, nil
}
