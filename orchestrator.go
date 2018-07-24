package model

import "errors"

type Orchestrator struct {
	root       *Environment
	Docker     attributes
	Parameters attributes
	Envvars    envvars
	Component
	Name string
}

func createOrchestrator(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	if yamlO.Repository == "" {
		vErrs.AddError(errors.New("no orchestrator specified"), "orchestrator")
		return Orchestrator{}
	} else {
		return Orchestrator{
			Component:  createComponent(vErrs, env.Lagoon, "orchestrator", yamlO.Repository, yamlO.Version),
			Docker:     createAttributes(yamlO.Docker, nil),
			Parameters: createAttributes(yamlO.Params, nil),
			Envvars:    createEnvvars(yamlO.Envvars, nil),
			root:       env}
	}
}
