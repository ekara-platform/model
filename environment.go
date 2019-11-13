package model

import (
	"strings"

	"gopkg.in/yaml.v2"
)

//go:generate go run ./generate/generate.go

type (
	//Environment represents an environment build based on a descriptor
	Environment struct {
		// The location of the environment root
		location DescriptorLocation `yaml:",omitempty"`
		// The environment name
		Name string `yaml:",omitempty"`
		// The environment qualifier
		Qualifier string `yaml:",omitempty"`
		// The environment description
		Description string `yaml:",omitempty"`
		// Ekara platform settings
		// The platform is "private" in order to be ignored
		// by the yaml unmarshalling process because an envirronment
		// cannot define component
		ekara *Platform
		// The descriptor variables
		Vars Parameters `yaml:",omitempty"`
		// The orchestrator used to manage the environment
		Orchestrator Orchestrator `yaml:",omitempty"`
		// The providers where to create the environment node sets
		Providers Providers `yaml:",omitempty"`
		// The node sets to create
		NodeSets NodeSets `yaml:",omitempty"`
		// The software stacks to install on the created node sets
		Stacks Stacks `yaml:",omitempty"`
		// The tasks which can be ran against the environment
		Tasks Tasks `yaml:",omitempty"`
		// The hooks linked to the environment lifecycle events
		Hooks EnvironmentHooks `yaml:",omitempty"`
		// The global volumes of the environment
		Volumes GlobalVolumes `yaml:",omitempty"`
		// Templates contains the templates defined into a descriptor
		Templates Patterns `yaml:",omitempty"`

		parcels []Parcel
	}

	//Parcel represent an environment intermediate version
	Parcel struct {
		ID    string
		Lines []string
	}
)

//CreateEnvironment creates a new environment based on the provided yaml
//The older passed as parameter y the name of the component holding the ekara.yaml on which
// the environment has been built
func CreateEnvironment(location string, yamlEnv yamlEnvironment, holder string) (*Environment, error) {
	env := &Environment{}
	var err error

	env.location = DescriptorLocation{Descriptor: location}
	env.Name = yamlEnv.Name
	env.Qualifier = yamlEnv.Qualifier
	env.Description = yamlEnv.Description
	env.Templates = createPatterns(env, env.location.appendPath("templates_patterns"), yamlEnv.Templates)
	env.Vars = CreateParameters(yamlEnv.yamlVars.Vars)

	env.Tasks, err = createTasks(env, env.location.appendPath("tasks"), &yamlEnv)
	if err != nil {
		return env, err
	}
	env.Orchestrator, err = createOrchestrator(env, env.location.appendPath("orchestrator"), &yamlEnv)
	if err != nil {
		return env, err
	}
	env.Providers, err = createProviders(env, env.location.appendPath("providers"), &yamlEnv)
	if err != nil {
		return env, err
	}
	env.NodeSets, err = createNodeSets(env, env.location.appendPath("nodes"), &yamlEnv)
	if err != nil {
		return env, err
	}
	// Only the main descriptor or a parent is allowed to define stacks
	if holder == MainComponentId || strings.HasPrefix(holder, EkaraComponentId) {
		env.Stacks, err = createStacks(env, holder, env.location.appendPath("stacks"), &yamlEnv)
		if err != nil {
			return env, err
		}
	}
	env.Hooks.Provision, err = createHook(env, env.location.appendPath("hooks.provision"), yamlEnv.Hooks.Provision)
	if err != nil {
		return env, err
	}
	env.Hooks.Deploy, err = createHook(env, env.location.appendPath("hooks.deploy"), yamlEnv.Hooks.Deploy)
	if err != nil {
		return env, err
	}
	env.Volumes = createGlobalVolumes(env, env.location.appendPath("volumes"), &yamlEnv)
	return env, nil

}

//Customize merges the content of the giver environment into the receiver
//
// Note: basic informations (name, qualifier, description) are only accepted once if the are not already defined
func (r *Environment) Customize(from Component, with *Environment) error {

	// We don't want to customize the templates defined into the environment
	// But instead we want to keep them into the component
	r.Platform().KeepTemplates(from, with.Templates)
	with.Templates = Patterns{}

	// basic informations (name, qualifier, description) are only accepted once if the are not already defined
	if r.Name == "" {
		r.Name = with.Name
	}

	if r.Qualifier == "" {
		r.Qualifier = with.Qualifier
	}

	if r.Description == "" {
		r.Description = with.Description
	}

	if err := r.Orchestrator.customize(with.Orchestrator); err != nil {
		return err
	}

	prs, err := r.Providers.customize(r, with.Providers)
	if err != nil {
		return err
	}
	r.Providers = prs

	nds, err := r.NodeSets.customize(r, with.NodeSets)
	if err != nil {
		return err
	}
	r.NodeSets = nds

	sts, err := r.Stacks.customize(r, with.Stacks)
	if err != nil {
		return err
	}
	r.Stacks = sts

	tas, err := r.Tasks.customize(r, with.Tasks)
	if err != nil {
		return err
	}
	r.Tasks = tas

	r.Vars = r.Vars.inherit(with.Vars)

	err = r.Hooks.customize(with.Hooks)

	l, err := lines(*r)
	if err != nil {
		return err
	}
	r.parcels = append(r.parcels, Parcel{ID: from.Id, Lines: l})

	return err
}

func lines(e Environment) ([]string, error) {
	yamlT, err := yaml.Marshal(e)
	if err != nil {
		return []string{}, err
	}
	ls := strings.Split(string(yamlT), "\n")
	return ls, nil
}

//Validate validate an environment
func (r Environment) Validate() ValidationErrors {
	vErrs := ValidationErrors{}

	vEr, e, _ := ErrorOnEmptyOrInvalid(r.Name, r.location.appendPath("name"), "empty environment name")
	vErrs.merge(vEr)
	if !e {
		vErrs.merge(ErrorOnInvalid(r.QualifiedName()))
	}

	vErrs.merge(ErrorOnInvalid(r.Orchestrator))

	vEr, _, _ = ErrorOnEmptyOrInvalid(r.Providers, r.location.appendPath("providers"), "no provider specified")
	vErrs.merge(vEr)

	vEr, _, _ = ErrorOnEmptyOrInvalid(r.NodeSets, r.location.appendPath("nodes"), "no node specified")
	vErrs.merge(vEr)

	vEr, _, _ = WarningOnEmptyOrInvalid(r.Stacks, r.location.appendPath("stacks"), "no stack specified")
	vErrs.merge(vEr)

	vErrs.merge(ErrorOnInvalid(r.Tasks))
	vErrs.merge(ErrorOnInvalid(r.Hooks))
	return vErrs
}

//InitEnvironment creates an new Environment
func InitEnvironment() *Environment {
	env := &Environment{
		ekara: &Platform{
			Components: make(map[string]Component),
		},
		parcels: make([]Parcel, 0, 0),
	}
	env.Orchestrator.cRef.env = env
	return env
}

//Platform Returns the platform on which the environment is built
func (r *Environment) Platform() *Platform {
	return r.ekara
}

//GetParcels returns environment intermediate versions.
//The last parcel, in the returned array, will be the final environment version
func (r *Environment) GetParcels() []Parcel {
	return r.parcels
}
