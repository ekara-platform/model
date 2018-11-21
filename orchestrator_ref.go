package model

type (
	OrchestratorRef struct {
		parameters Parameters
		docker     Parameters
		envVars    EnvVars

		env      *Environment
		location DescriptorLocation
	}
)

func (r OrchestratorRef) validate() ValidationErrors {
	return ValidationErrors{}
}

func (r *OrchestratorRef) merge(other OrchestratorRef) error {
	return nil
}

func (r OrchestratorRef) Resolve() (Orchestrator, error) {
	validationErrors := r.validate()
	if validationErrors.HasErrors() {
		return Orchestrator{}, validationErrors
	}
	orchestrator := r.env.Orchestrator
	return Orchestrator{
		Component:  orchestrator.Component,
		Parameters: r.parameters.inherits(orchestrator.Parameters),
		Docker:     r.docker.inherits(orchestrator.Docker),
		EnvVars:    r.envVars.inherits(orchestrator.EnvVars)}, nil
}

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) OrchestratorRef {
	return OrchestratorRef{
		env:        env,
		parameters: createParameters(yamlRef.Params),
		docker:     createParameters(yamlRef.Docker),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
	}
}
