package model

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"

	"github.com/imdario/mergo"
)

type Environment struct {
	Component

	// The environment name
	Name string
	// The environment qualifier
	Qualifier string

	// The environment description
	Description string

	// Ekara platform settings
	Ekara EkaraPlatform

	// The providers where to create the environment node sets
	Providers map[string]Provider
	// The orchestrator used to manage the environment
	Orchestrator Orchestrator
	// The nodesets to create
	NodeSets map[string]NodeSet
	// The software stacks to install on the created node sets
	Stacks map[string]Stack
	// The tasks which can be ran against the environment
	Tasks map[string]Task

	Hooks EnvironmentHooks
}

func (env *Environment) Merge(other *Environment) {
	mergo.Merge(env, other, mergo.WithOverride)
}

func (r Environment) MarshalJSON() ([]byte, error) {
	t := struct {
		Name          string         `json:",omitempty"`
		Qualifier     string         `json:",omitempty"`
		QualifiedName string         `json:",omitempty"`
		Description   string         `json:",omitempty"`
		Ekara         *EkaraPlatform `json:",omitempty"`
		Providers     map[string]Provider
		Orchestrator  *Orchestrator `json:",omitempty"`
		NodeSets      map[string]NodeSet
		Stacks        map[string]Stack
		Tasks         map[string]Task
		Hooks         *EnvironmentHooks `json:",omitempty"`
	}{
		Name:          r.Name,
		Qualifier:     r.Qualifier,
		QualifiedName: r.QualifiedName().String(),
		Description:   r.Description,

		Ekara:        &r.Ekara,
		Providers:    r.Providers,
		Orchestrator: &r.Orchestrator,
		NodeSets:     r.NodeSets,
		Stacks:       r.Stacks,
		Tasks:        r.Tasks,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

type EnvironmentHooks struct {
	Init      Hook
	Provision Hook
	Deploy    Hook
	Undeploy  Hook
	Destroy   Hook
}

func (r EnvironmentHooks) HasTasks() bool {
	return r.Init.HasTasks() ||
		r.Provision.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks() ||
		r.Destroy.HasTasks()
}

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

func Parse(logger *log.Logger, u *url.URL) (*Environment, error) {
	return ParseWithData(logger, u, map[string]interface{}{})
}

func ParseWithData(logger *log.Logger, u *url.URL, data map[string]interface{}) (*Environment, error) {
	vErrs := ValidationErrors{}
	if hasSuffixIgnoringCase(u.Path, ".yaml") || hasSuffixIgnoringCase(u.Path, ".yml") {
		var yamlEnv yamlEnvironment
		yamlEnv, err := parseYamlDescriptor(logger, u, data)
		if err != nil {
			return &Environment{}, err
		}
		env := createEnvironment(&vErrs, &yamlEnv)
		if vErrs.HasErrors() || vErrs.HasWarnings() {
			return &env, vErrs
		} else {
			return &env, nil
		}
	} else {
		return &Environment{}, errors.New("unsupported file format")
	}
}

func createEnvironment(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Environment {
	var env = Environment{}
	if len(yamlEnv.Name) == 0 {
		vErrs.AddError(errors.New("empty environment name"), "name")
	}

	// Check the validity of the qualified name of the environment
	if !yamlEnv.QualifiedName().ValidQualifiedName() {
		vErrs.AddError(errors.New("The environment name or the qualifier contain a non alphanumeric character"), "name_qualifier")
	}

	env.Name = yamlEnv.Name
	env.Qualifier = yamlEnv.Qualifier
	env.Description = yamlEnv.Description

	env.Ekara = createEkaraPlatform(vErrs, yamlEnv)
	env.Tasks = createTasks(vErrs, &env, yamlEnv)
	env.Orchestrator = createOrchestrator(vErrs, &env, yamlEnv)
	env.Providers = createProviders(vErrs, &env, yamlEnv)
	env.NodeSets = createNodeSets(vErrs, &env, yamlEnv)
	env.Stacks = createStacks(vErrs, &env, yamlEnv)
	env.Hooks.Init = createHook(vErrs, "hooks.init", &env, yamlEnv.Hooks.Init)
	env.Hooks.Provision = createHook(vErrs, "hooks.provision", &env, yamlEnv.Hooks.Provision)
	env.Hooks.Deploy = createHook(vErrs, "hooks.deploy", &env, yamlEnv.Hooks.Deploy)
	env.Hooks.Undeploy = createHook(vErrs, "hooks.undeploy", &env, yamlEnv.Hooks.Undeploy)
	env.Hooks.Destroy = createHook(vErrs, "hooks.destroy", &env, yamlEnv.Hooks.Destroy)
	return env
}
