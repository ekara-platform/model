package model

type (
	//Environment represents an environment build based on a descriptor
	Environment struct {
		OriginalEnv yamlEnvironment
		// The location of the environment root
		location DescriptorLocation `yaml:",omitempty"`
		// The environment name
		Name string
		// The environment qualifier
		Qualifier string
		// The environment description
		Description string
		// Ekara platform settings
		Ekara *Platform
		// The descriptor variables
		Vars Parameters
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
		// The hooks linked to the environment lifecycle events
		Hooks EnvironmentHooks
		// The global volumes of the environment
		Volumes GlobalVolumes
		// Templates contains the templates defined into a descriptor
		Templates Patterns
	}
)

//CreateEnvironment creates a new environment
//	Parameters
//
//		url: 	The complete url pointing on the descritor used to build the environment.
//			The two only supported extension are ".yaml" and ".yml"!
//		data: The context used to substitute variables into the environment descriptor
//
func CreateEnvironment(url EkUrl, data *TemplateContext) (*Environment, error) {
	env := &Environment{}
	var err error
	var yamlEnv yamlEnvironment

	yamlEnv, err = parseYamlDescriptor(url, data)
	if err != nil {
		return env, err
	}
	env.OriginalEnv = yamlEnv

	env.location = DescriptorLocation{Descriptor: url.String()}
	env.Name = yamlEnv.Name
	env.Qualifier = yamlEnv.Qualifier
	env.Description = yamlEnv.Description
	env.Templates = createPatterns(env, env.location.appendPath("templates_patterns"), yamlEnv.Templates)
	env.Ekara, err = createPlatform(&yamlEnv)
	if err != nil {
		return env, err
	}

	vars, err := CreateParameters(yamlEnv.yamlVars.Vars)
	if err != nil {
		return env, err
	}
	env.Vars = vars

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
	env.Stacks, err = createStacks(env, env.location.appendPath("stacks"), &yamlEnv)
	if err != nil {
		return env, err
	}
	env.Hooks.Provision, err = createHook(env, env.location.appendPath("hooks.provision"), yamlEnv.Hooks.Provision)
	if err != nil {
		return env, err
	}
	env.Hooks.Deploy, err = createHook(env, env.location.appendPath("hooks.deploy"), yamlEnv.Hooks.Deploy)
	if err != nil {
		return env, err
	}
	env.Hooks.Undeploy, err = createHook(env, env.location.appendPath("hooks.undeploy"), yamlEnv.Hooks.Undeploy)
	if err != nil {
		return env, err
	}
	env.Hooks.Destroy, err = createHook(env, env.location.appendPath("hooks.destroy"), yamlEnv.Hooks.Destroy)
	env.Volumes = createGlobalVolumes(env, env.location.appendPath("volumes"), &yamlEnv)
	return env, nil

}

//Merge merges the content of the other environment into the receiver
//
// Note: basic informations (name, qualifier, description) are only accepted in root descriptor
func (r *Environment) Merge(other *Environment) error {
	// basic informations (name, qualifier, description) are only accepted in root descriptor
	if err := r.Ekara.merge(*other.Ekara); err != nil {
		return err
	}

	if err := r.Orchestrator.merge(other.Orchestrator); err != nil {
		return err
	}

	if prs, err := r.Providers.merge(r, other.Providers); err != nil {
		return err
	} else {
		r.Providers = prs
	}

	if nds, err := r.NodeSets.merge(r, other.NodeSets); err != nil {
		return err
	} else {
		r.NodeSets = nds
	}
	if sts, err := r.Stacks.merge(r, other.Stacks); err != nil {
		return err
	} else {
		r.Stacks = sts
	}
	if tas, err := r.Tasks.merge(r, other.Tasks); err != nil {
		return err
	} else {
		r.Tasks = tas
	}
	return r.Hooks.merge(other.Hooks)
}

//Validate validate an environment
func (r Environment) Validate() ValidationErrors {
	vErrs := ValidationErrors{}

	vEr, e, _ := ErrorOnEmptyOrInvalid(r.Name, r.location.appendPath("name"), "empty environment name")
	vErrs.merge(vEr)
	if !e {
		vErrs.merge(ErrorOnInvalid(r.QualifiedName()))
	}

	vErrs.merge(ErrorOnInvalid(r.Ekara))
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
