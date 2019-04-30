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

func (r orchestratorRef) Resolve() (Orchestrator, error) {
	orchestrator := r.env.Orchestrator
	return Orchestrator{
		Component:  orchestrator.Component,
		Parameters: r.parameters.inherits(orchestrator.Parameters),
		Docker:     r.docker.inherits(orchestrator.Docker),
		EnvVars:    r.envVars.inherits(orchestrator.EnvVars)}, nil
}

func (r *orchestratorRef) merge(other orchestratorRef) error {
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
	r.docker = r.docker.inherits(other.docker)
	return nil
}

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) orchestratorRef {
	return orchestratorRef{
		env:        env,
		parameters: createParameters(yamlRef.Params),
		docker:     createParameters(yamlRef.Docker),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
	}
}

// OrchestratorParams returns the parameters required to install the orchestrator
func (r orchestratorRef) OrchestratorParams() (map[string]interface{}, error) {
	op := make(map[string]interface{})
	o, err := r.Resolve()
	if err != nil {
		return op, err
	}

	op["docker"] = o.Docker
	op["params"] = o.Parameters
	return op, nil
}
