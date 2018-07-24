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
	o := Orchestrator{}
	yamlO := yamlEnv.Orchestrator
	if yamlO.Name == "" {
		vErrs.AddError(errors.New("no orchestrator specified"), "orchestrator")
	} else {
		o.Component = createComponent(vErrs, env, "orchestrator", yamlO.Repository, yamlO.Version)
		o.Name = yamlO.Name
		o.Docker = createAttributes(yamlO.Docker, nil)
		o.Parameters = createAttributes(yamlO.Params, nil)
		o.Envvars = createEnvvars(yamlO.Envvars, nil)
		o.root = env
	}
	return o
}
