package model

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
)

type (
	Environment struct {
		// Validation errors that occurred during the building of the environment
		errors ValidationErrors
		// The location of the environment root
		location DescriptorLocation

		// Global imports
		Imports []string
		// The environment name
		Name string
		// The environment qualifier
		Qualifier string
		// The environment description
		Description string
		// Ekara platform settings
		Ekara Platform
		// The orchestrator used to manage the environment
		Orchestrator Orchestrator
		// The providers where to create the environment node sets
		Providers Providers
		// The node sets to create
		NodeSets NodeSets
		// The software stacks to install on the created node sets
		Stacks Stacks
		// The tasks which can be ran against the environment
		Tasks Tasks
		// Global environment hooks
		Hooks EnvironmentHooks
	}

	EnvironmentHooks struct {
		Init      Hook
		Provision Hook
		Deploy    Hook
		Undeploy  Hook
		Destroy   Hook
	}
)

func (r Environment) MarshalJSON() ([]byte, error) {
	t := struct {
		Name          string    `json:",omitempty"`
		Qualifier     string    `json:",omitempty"`
		QualifiedName string    `json:",omitempty"`
		Description   string    `json:",omitempty"`
		Ekara         *Platform `json:",omitempty"`
		Providers     *Providers
		Orchestrator  *Orchestrator `json:",omitempty"`
		NodeSets      *NodeSets
		Stacks        *Stacks
		Tasks         *Tasks
		Hooks         *EnvironmentHooks `json:",omitempty"`
	}{
		Name:          r.Name,
		Qualifier:     r.Qualifier,
		QualifiedName: r.QualifiedName().String(),
		Description:   r.Description,

		Ekara:        &r.Ekara,
		Providers:    &r.Providers,
		Orchestrator: &r.Orchestrator,
		NodeSets:     &r.NodeSets,
		Stacks:       &r.Stacks,
		Tasks:        &r.Tasks,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

func (r EnvironmentHooks) HasTasks() bool {
	return r.Init.HasTasks() ||
		r.Provision.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r *EnvironmentHooks) merge(other EnvironmentHooks) error {
	if err := r.Init.merge(other.Init); err != nil {
		return err
	}
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	if err := r.Deploy.merge(other.Deploy); err != nil {
		return err
	}
	if err := r.Undeploy.merge(other.Undeploy); err != nil {
		return err
	}
	if err := r.Destroy.merge(other.Destroy); err != nil {
		return err
	}
	return nil
}

func (r EnvironmentHooks) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(r.Init.validate())
	vErrs.merge(r.Provision.validate())
	vErrs.merge(r.Deploy.validate())
	vErrs.merge(r.Undeploy.validate())
	vErrs.merge(r.Destroy.validate())
	return vErrs
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

func CreateEnvironment(logger *log.Logger, url *url.URL, data map[string]interface{}) (Environment, error) {
	env := Environment{}
	if hasSuffixIgnoringCase(url.Path, ".yaml") || hasSuffixIgnoringCase(url.Path, ".yml") {
		var yamlEnv yamlEnvironment
		yamlEnv, err := parseYamlDescriptor(logger, url, data)
		if err != nil {
			return env, err
		}
		env.location = DescriptorLocation{Descriptor: url.String()}
		env.Imports = yamlEnv.Imports
		env.Name = yamlEnv.Name
		env.Qualifier = yamlEnv.Qualifier
		env.Description = yamlEnv.Description
		env.Ekara = createPlatform(&env, &yamlEnv)
		env.Tasks = createTasks(&env, &yamlEnv)
		env.Orchestrator = createOrchestrator(&env, &yamlEnv)
		env.Providers = createProviders(&env, &yamlEnv)
		env.NodeSets = createNodeSets(&env, &yamlEnv)
		env.Stacks = createStacks(&env, &yamlEnv)
		env.Hooks.Init = createHook(&env, env.location.appendPath("hooks.init"), yamlEnv.Hooks.Init)
		env.Hooks.Provision = createHook(&env, env.location.appendPath("hooks.provision"), yamlEnv.Hooks.Provision)
		env.Hooks.Deploy = createHook(&env, env.location.appendPath("hooks.deploy"), yamlEnv.Hooks.Deploy)
		env.Hooks.Undeploy = createHook(&env, env.location.appendPath("hooks.undeploy"), yamlEnv.Hooks.Undeploy)
		env.Hooks.Destroy = createHook(&env, env.location.appendPath("hooks.destroy"), yamlEnv.Hooks.Destroy)
	} else {
		return env, errors.New("unsupported file format")
	}
	return env, nil
}

func (r *Environment) Merge(other Environment) error {
	// Data and basic info (name, qualifier, description) are only accepted in root descriptor
	if err := r.Ekara.merge(other.Ekara); err != nil {
		return err
	}
	if err := r.Orchestrator.merge(other.Orchestrator); err != nil {
		return err
	}
	if err := r.Providers.merge(r, other.Providers); err != nil {
		return err
	}
	if err := r.NodeSets.merge(r, other.NodeSets); err != nil {
		return err
	}
	if err := r.Stacks.merge(r, other.Stacks); err != nil {
		return err
	}
	if err := r.Tasks.merge(r, other.Tasks); err != nil {
		return err
	}
	if err := r.Hooks.merge(other.Hooks); err != nil {
		return err
	}
	return nil
}

func (r Environment) Validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(r.errors)
	vEr, e, _ := ErrorOnEmpty(r.Name, r.location.appendPath("name"), "empty environment name")
	vErrs.merge(vEr)
	if !e && !r.QualifiedName().ValidQualifiedName() {
		vErrs.addError(errors.New("the environment name or the qualifier contains a non alphanumeric character"), r.location.appendPath("name|qualifier"))
	}
	vErrs.merge(r.Ekara.validate())
	vErrs.merge(r.Orchestrator.validate())

	vEr, _, _ = ErrorOnEmpty(r.Providers, r.location.appendPath("providers"), "no provider specified")
	vErrs.merge(vEr)
	vEr, _, _ = ErrorOnEmpty(r.NodeSets, r.location.appendPath("nodes"), "no node specified")
	vErrs.merge(vEr)
	vEr, _, _ = WarningOnEmpty(r.Stacks, r.location.appendPath("stacks"), "no stack specified")
	vErrs.merge(vEr)

	for _, t := range r.Tasks {
		vErrs.merge(t.validate())
	}
	vErrs.merge(r.Hooks.validate())
	return vErrs
}
