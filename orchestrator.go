package model

import (
	"encoding/json"
)

type Orchestrator struct {
	// The component containing the orchestrator
	Component ComponentRef
	// The orchestrator parameters
	Parameters Parameters
	// The Docker parameters
	Docker Parameters
	// The orchestrator environment variables
	EnvVars EnvVars
}

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

func (r *Orchestrator) merge(other Orchestrator) {
	r.Component.merge(other.Component)
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.Docker = r.Docker.inherits(other.Docker)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
}

func (r Orchestrator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Component  string     `json:",omitempty"`
		Parameters Parameters `json:",omitempty"`
		Docker     Parameters `json:",omitempty"`
		EnvVars    EnvVars    `json:",omitempty"`
	}{
		Component:  r.Component.Resolve().Id,
		Parameters: r.Parameters,
		Docker:     r.Docker,
		EnvVars:    r.EnvVars,
	})
}

type OrchestratorRef struct {
	parameters Parameters
	docker     Parameters
	envVars    EnvVars

	env      *Environment
	location DescriptorLocation
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

func (r *OrchestratorRef) merge(other OrchestratorRef) {

}

func (r OrchestratorRef) Resolve() Orchestrator {
	validationErrors := r.validate()
	if validationErrors.HasErrors() {
		panic(validationErrors)
	}
	orchestrator := r.env.Orchestrator
	return Orchestrator{
		Component:  orchestrator.Component,
		Parameters: r.parameters.inherits(orchestrator.Parameters),
		Docker:     r.docker.inherits(orchestrator.Docker),
		EnvVars:    r.envVars.inherits(orchestrator.EnvVars)}
}

// OrchestratorParams returns the parameters required to install the orchestrator
func (r OrchestratorRef) OrchestratorParams() map[string]interface{} {
	o := r.Resolve()
	op := make(map[string]interface{})
	op["docker"] = o.Docker
	op["params"] = o.Parameters
	return op
}
