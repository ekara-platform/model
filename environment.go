package model

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
)

type Reference interface {
	attach(env *Environment)
}

type Environment struct {
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
	Providers map[string]Provider
	// The node sets to create
	NodeSets map[string]NodeSet
	// The software stacks to install on the created node sets
	Stacks map[string]Stack
	// The tasks which can be ran against the environment
	Tasks map[string]Task
	// Global hooks
	Hooks EnvironmentHooks
}

func (r Environment) MarshalJSON() ([]byte, error) {
	t := struct {
		Name          string    `json:",omitempty"`
		Qualifier     string    `json:",omitempty"`
		QualifiedName string    `json:",omitempty"`
		Description   string    `json:",omitempty"`
		Ekara         *Platform `json:",omitempty"`
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

func CreateEnvironment(logger *log.Logger, u *url.URL, data map[string]interface{}) (Environment, error) {
	env := Environment{}
	err := env.parse(logger, u, data)
	if err != nil {
		return env, err
	}
	return env, nil
}

func (r *Environment) Merge(other Environment) error {
	// Data and basic info (name, qualifier, description) are only accepted in root descriptor
	r.Ekara.merge(other.Ekara)
	r.Orchestrator.merge(other.Orchestrator)
	for id, p := range other.Providers {
		if provider, ok := r.Providers[id]; ok {
			provider.merge(p)
		} else {
			p.Component.env = r
			r.Providers[id] = p
		}
	}
	for id, n := range other.NodeSets {
		if nodeSet, ok := r.NodeSets[id]; ok {
			nodeSet.merge(n)
		} else {
			n.Provider.env = r
			n.Orchestrator.env = r
			r.NodeSets[id] = n
		}
	}
	for id, s := range other.Stacks {
		if stack, ok := r.Stacks[id]; ok {
			stack.merge(s)
		} else {
			s.Component.env = r
			r.Stacks[id] = s
		}
	}
	for id, t := range other.Tasks {
		if task, ok := r.Tasks[id]; ok {
			task.merge(t)
		} else {
			t.Component.env = r
			r.Tasks[id] = t
		}
	}
	r.Hooks.Init.merge(r.Hooks.Init)
	r.Hooks.Provision.merge(r.Hooks.Provision)
	r.Hooks.Deploy.merge(r.Hooks.Deploy)
	r.Hooks.Undeploy.merge(r.Hooks.Undeploy)
	r.Hooks.Destroy.merge(r.Hooks.Destroy)
	return nil
}

func (r Environment) Validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(r.errors)
	if len(r.Name) == 0 {
		vErrs.addError(errors.New("empty environment name"), r.location.appendPath("name"))
	}
	if !r.QualifiedName().ValidQualifiedName() {
		vErrs.addError(errors.New("the environment name or the qualifier contains a non alphanumeric character"), r.location.appendPath("name|qualifier"))
	}
	vErrs.merge(r.Ekara.validate())
	vErrs.merge(r.Orchestrator.validate())
	if len(r.Providers) == 0 {
		vErrs.addError(errors.New("no provider specified"), r.location.appendPath("providers"))
	}
	for _, p := range r.Providers {
		vErrs.merge(p.validate())
	}
	if len(r.NodeSets) == 0 {
		vErrs.addError(errors.New("no node specified"), r.location.appendPath("nodes"))
	}
	for _, n := range r.NodeSets {
		vErrs.merge(n.validate())
	}
	if len(r.Stacks) == 0 {
		vErrs.addWarning("no stack specified", r.location.appendPath("stacks"))
	}
	for _, s := range r.Stacks {
		vErrs.merge(s.validate())
	}
	for _, t := range r.Tasks {
		vErrs.merge(t.validate())
	}
	vErrs.merge(r.Hooks.Init.validate())
	vErrs.merge(r.Hooks.Provision.validate())
	vErrs.merge(r.Hooks.Deploy.validate())
	vErrs.merge(r.Hooks.Undeploy.validate())
	vErrs.merge(r.Hooks.Destroy.validate())
	return vErrs
}

func (r *Environment) parse(logger *log.Logger, u *url.URL, data map[string]interface{}) error {
	if hasSuffixIgnoringCase(u.Path, ".yaml") || hasSuffixIgnoringCase(u.Path, ".yml") {
		var yamlEnv yamlEnvironment
		yamlEnv, err := parseYamlDescriptor(logger, u, data)
		if err != nil {
			return err
		}
		r.build(u, &yamlEnv)
		return nil
	} else {
		return errors.New("unsupported file format")
	}
}

func (r *Environment) build(url *url.URL, yamlEnv *yamlEnvironment) {
	r.location = DescriptorLocation{Path: "", Descriptor: url.String()}
	r.Imports = yamlEnv.Imports
	r.Name = yamlEnv.Name
	r.Qualifier = yamlEnv.Qualifier
	r.Description = yamlEnv.Description
	r.Ekara = createPlatform(r, yamlEnv)
	r.Tasks = createTasks(r, yamlEnv)
	r.Orchestrator = createOrchestrator(r, yamlEnv)
	r.Providers = createProviders(r, yamlEnv)
	r.NodeSets = createNodeSets(r, yamlEnv)
	r.Stacks = createStacks(r, yamlEnv)
	r.Hooks.Init = createHook(r, r.location.appendPath("hooks.init"), yamlEnv.Hooks.Init)
	r.Hooks.Provision = createHook(r, r.location.appendPath("hooks.provision"), yamlEnv.Hooks.Provision)
	r.Hooks.Deploy = createHook(r, r.location.appendPath("hooks.deploy"), yamlEnv.Hooks.Deploy)
	r.Hooks.Undeploy = createHook(r, r.location.appendPath("hooks.undeploy"), yamlEnv.Hooks.Undeploy)
	r.Hooks.Destroy = createHook(r, r.location.appendPath("hooks.destroy"), yamlEnv.Hooks.Destroy)
}
