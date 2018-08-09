package model

import "errors"

type Orchestrator struct {
	// The component containing the orchestrator
	Component ComponentRef
	// The orchestrator parameters
	Parameters Parameters
	// The orchestrator environment variables
	EnvVars EnvVars
}

type OrchestratorRef struct {
	orchestrator *Orchestrator
	parameters   Parameters
	envVars      EnvVars
}

func (o OrchestratorRef) Resolve() Orchestrator {
	return Orchestrator{
		Component:  o.orchestrator.Component,
		Parameters: o.parameters.inherit(o.orchestrator.Parameters),
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
			EnvVars:    createEnvVars(yamlO.Env)}
	}
}

func createOrchestratorRef(env *Environment, yamlRef yamlOrchestratorRef) OrchestratorRef {
	return OrchestratorRef{
		orchestrator: &env.Orchestrator,
		parameters:   createParameters(yamlRef.Params),
		envVars:      createEnvVars(yamlRef.Env)}
}
