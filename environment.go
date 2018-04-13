package model

import (
	"log"
	"strings"
	"errors"
)

type Environment struct {
	Labels
	Component

	// Global attributes
	Name        string
	Description string
	Version     Version

	// Lagoon platform attributes
	Lagoon Lagoon

	// Definition attributes
	Providers map[string]Provider
	NodeSets  map[string]NodeSet
	Stacks    map[string]Stack
	Tasks     map[string]Task

	Hooks struct {
		Init      Hook
		Provision Hook
		Deploy    Hook
		Undeploy  Hook
		Destroy   Hook
	}
}

func Parse(logger *log.Logger, location string) (env Environment, err error, vErrs ValidationErrors) {
	vErrs = ValidationErrors{}
	if strings.HasSuffix(strings.ToUpper(location), ".YAML") ||
		strings.HasSuffix(strings.ToUpper(location), ".YML") {
		var yamlEnv yamlEnvironment
		yamlEnv, err = parseYamlDescriptor(logger, location)
		if err != nil {
			return
		}
		env = createEnvironment(&vErrs, &yamlEnv)
		postValidate(&vErrs, &env)
		if vErrs.HasErrors() {
			err = errors.New("validation errors have occurred")
		}
	} else {
		err = errors.New("unsupported file format")
	}
	return
}

func createEnvironment(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Environment {
	var env = Environment{}

	env.Name = yamlEnv.Name
	env.Description = yamlEnv.Description
	env.Labels = createLabels(vErrs, yamlEnv.Labels...)
	env.Version = createVersion(vErrs, "version", yamlEnv.Version)
	env.Lagoon = createLagoon(vErrs, &env, yamlEnv)
	env.Tasks = createTasks(vErrs, &env, yamlEnv)
	env.Providers = createProviders(vErrs, &env, yamlEnv)
	env.NodeSets = createNodeSets(vErrs, &env, yamlEnv)
	env.Stacks = createStacks(vErrs, &env, yamlEnv)
	env.Hooks.Init = createHook(vErrs, env.Tasks, "hooks.init", yamlEnv.Hooks.Init)
	env.Hooks.Provision = createHook(vErrs, env.Tasks, "hooks.provision", yamlEnv.Hooks.Provision)
	env.Hooks.Deploy = createHook(vErrs, env.Tasks, "hooks.deploy", yamlEnv.Hooks.Deploy)
	env.Hooks.Undeploy = createHook(vErrs, env.Tasks, "hooks.undeploy", yamlEnv.Hooks.Undeploy)
	env.Hooks.Destroy = createHook(vErrs, env.Tasks, "hooks.destroy", yamlEnv.Hooks.Destroy)
	return env
}

func postValidate(vErrs *ValidationErrors, env *Environment) ValidationErrors {
	validationErrors := ValidationErrors{}
	return validationErrors
}
