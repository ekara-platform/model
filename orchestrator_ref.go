package model

type (
	//OrchestratorRef represents a reference on an Orchestrator
	OrchestratorRef struct {
		parameters Parameters
		envVars    EnvVars
		env        *Environment
		location   DescriptorLocation
	}
)

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) (OrchestratorRef, error) {
	return OrchestratorRef{
		env:        env,
		parameters: CreateParameters(yamlRef.Params),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
	}, nil
}

func (r *OrchestratorRef) customize(with OrchestratorRef) error {
	r.parameters = with.parameters.inherit(r.parameters)
	r.envVars = with.envVars.inherit(r.envVars)
	return nil
}

//Resolve returns the referenced Orchestrator
func (r OrchestratorRef) Resolve() (Orchestrator, error) {
	orchestrator := r.env.Orchestrator
	params := r.parameters.inherit(orchestrator.Parameters)
	envVars := r.envVars.inherit(orchestrator.EnvVars)
	return Orchestrator{
		cRef:       orchestrator.cRef,
		Parameters: params,
		EnvVars:    envVars,
	}, nil
}
