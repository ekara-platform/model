package model

type (
	orchestratorRef struct {
		parameters Parameters
		docker     Parameters
		envVars    EnvVars

		env      *Environment
		location DescriptorLocation
	}
)

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) orchestratorRef {
	return orchestratorRef{
		env:        env,
		parameters: createParameters(yamlRef.Params),
		docker:     createParameters(yamlRef.Docker),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
	}
}

func (r *orchestratorRef) merge(other orchestratorRef) error {
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
	r.docker = r.docker.inherits(other.docker)
	return nil
}

func (r orchestratorRef) Resolve() (Orchestrator, error) {
	orchestrator := r.env.Orchestrator
	return Orchestrator{
		cRef:       orchestrator.cRef,
		Parameters: r.parameters.inherits(orchestrator.Parameters),
		Docker:     r.docker.inherits(orchestrator.Docker),
		EnvVars:    r.envVars.inherits(orchestrator.EnvVars)}, nil
}
