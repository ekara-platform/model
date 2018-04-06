package descriptor

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

func ParseDescriptor(logger *log.Logger, location string) (env Environment, err error) {
	if strings.HasSuffix(strings.ToUpper(location), ".YAML") ||
		strings.HasSuffix(strings.ToUpper(location), ".YML") {
		var yamlEnv yamlEnvironment
		yamlEnv, err = parseYamlDescriptor(logger, location)
		if err != nil {
			return
		}
		env, err = createEnvironment(&yamlEnv)
	} else {
		err = errors.New("unsupported file format")
	}
	return
}

func createEnvironment(yamlEnv *yamlEnvironment) (env Environment, err error) {
	env = Environment{}

	env.Name = yamlEnv.Name
	env.Description = yamlEnv.Description
	env.Labels = createLabels(yamlEnv.Labels...)
	env.Version, err = createVersion(yamlEnv.Version)
	if err != nil {
		return
	}
	env.Lagoon, err = createLagoon(&env, yamlEnv)
	if err != nil {
		return
	}
	env.Tasks, err = createTasks(&env, yamlEnv)
	if err != nil {
		return
	}
	env.Providers, err = createProviders(&env, yamlEnv)
	if err != nil {
		return
	}
	env.NodeSets, err = createNodeSets(&env, yamlEnv)
	if err != nil {
		return
	}
	env.Stacks, err = createStacks(&env, yamlEnv)
	if err != nil {
		return
	}
	env.Hooks.Init, err = createHook(env.Tasks, yamlEnv.Hooks.Init)
	if err != nil {
		return
	}
	env.Hooks.Provision, err = createHook(env.Tasks, yamlEnv.Hooks.Provision)
	if err != nil {
		return
	}
	env.Hooks.Deploy, err = createHook(env.Tasks, yamlEnv.Hooks.Deploy)
	if err != nil {
		return
	}
	env.Hooks.Undeploy, err = createHook(env.Tasks, yamlEnv.Hooks.Undeploy)
	if err != nil {
		return
	}
	env.Hooks.Destroy, err = createHook(env.Tasks, yamlEnv.Hooks.Destroy)
	if err != nil {
		return
	}
	return
}
