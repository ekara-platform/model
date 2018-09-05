package model

import (
	"encoding/json"
	"errors"
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

func (r Orchestrator) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Component  *ComponentRef `json:",omitempty"`
		Parameters *Parameters   `json:",omitempty"`
		Docker     *Parameters   `json:",omitempty"`
		EnvVars    *EnvVars      `json:",omitempty"`
	}{
		Component:  &r.Component,
		Parameters: &r.Parameters,
		Docker:     &r.Docker,
		EnvVars:    &r.EnvVars,
	})
}

type OrchestratorRef struct {
	orchestrator *Orchestrator
	parameters   Parameters
	docker       Parameters
	envVars      EnvVars
}

// OrchestratorParams returns the parameters required to install the orchestrator
func (ref OrchestratorRef) OrchestratorParams() map[string]interface{} {
	o := ref.Resolve()
	r := make(map[string]interface{})
	r["docker"] = o.Docker
	r["params"] = o.Parameters
	return r
}

func (o OrchestratorRef) Resolve() Orchestrator {
	return Orchestrator{
		Component:  o.orchestrator.Component,
		Parameters: o.parameters.inherit(o.orchestrator.Parameters),
		Docker:     o.docker.inherit(o.orchestrator.Docker),
		EnvVars:    o.envVars.inherit(o.orchestrator.EnvVars)}
}

func createOrchestrator(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	if yamlO.Component == "" {
		vErrs.AddError(errors.New("no orchestrator specified"), "orchestrator")
		return Orchestrator{}
	} else {
		return Orchestrator{
			Component:  createComponentRef(vErrs, env.Lagoon.Components, "orchestrator", yamlO.Component),
			Parameters: createParameters(yamlO.Params),
			Docker:     createParameters(yamlO.Docker),
			EnvVars:    createEnvVars(yamlO.Env)}
	}
}

func createOrchestratorRef(env *Environment, yamlRef yamlOrchestratorRef) OrchestratorRef {
	return OrchestratorRef{
		orchestrator: &env.Orchestrator,
		parameters:   createParameters(yamlRef.Params),
		docker:       createParameters(yamlRef.Docker),
		envVars:      createEnvVars(yamlRef.Env)}
}
