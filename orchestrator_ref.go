package model

type (
	orchestratorRef struct {
		parameters Parameters
		docker     Parameters
		envVars    EnvVars

		env       *Environment
		location  DescriptorLocation
		templates Patterns
	}
)

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) (orchestratorRef, error) {
	oParams, err := CreateParameters(yamlRef.Params)
	if err != nil {
		return orchestratorRef{}, err
	}
	dParams, err := CreateParameters(yamlRef.Docker)
	if err != nil {
		return orchestratorRef{}, err
	}
	envVars, err := createEnvVars(yamlRef.Env)
	if err != nil {
		return orchestratorRef{}, err
	}
	return orchestratorRef{
		env:        env,
		parameters: oParams,
		docker:     dParams,
		envVars:    envVars,
		location:   location,
		templates:  createPatterns(env, location.appendPath("templates_patterns"), yamlRef.Templates),
	}, nil
}

func (r *orchestratorRef) merge(other orchestratorRef) error {
	var err error
	r.parameters, err = r.parameters.inherit(other.parameters)
	if err != nil {
		return err
	}
	r.envVars, err = r.envVars.inherit(other.envVars)
	if err != nil {
		return err
	}
	r.docker, err = r.docker.inherit(other.docker)
	if err != nil {
		return err
	}
	r.templates = r.templates.inherit(other.templates)
	return nil
}

func (r orchestratorRef) Resolve() (Orchestrator, error) {
	orchestrator := r.env.Orchestrator
	params, err := r.parameters.inherit(orchestrator.Parameters)
	if err != nil {
		return Orchestrator{}, err
	}
	docker, err := r.docker.inherit(orchestrator.Docker)
	if err != nil {
		return Orchestrator{}, err
	}
	envVars, err := r.envVars.inherit(orchestrator.EnvVars)
	if err != nil {
		return Orchestrator{}, err
	}
	return Orchestrator{
		cRef:       orchestrator.cRef,
		Parameters: params,
		Docker:     docker,
		EnvVars:    envVars,
		Templates:  r.templates.inherit(orchestrator.Templates),
	}, nil
}
