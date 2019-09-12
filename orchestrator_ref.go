package model

type (
	//OrchestratorRef represents a reference on an Orchestrator
	OrchestratorRef struct {
		parameters Parameters
		envVars    EnvVars

		env      *Environment
		location DescriptorLocation
	}
)

func createOrchestratorRef(env *Environment, location DescriptorLocation, yamlRef yamlOrchestratorRef) (OrchestratorRef, error) {
	oParams, err := CreateParameters(yamlRef.Params)
	if err != nil {
		return OrchestratorRef{}, err
	}
	envVars, err := createEnvVars(yamlRef.Env)
	if err != nil {
		return OrchestratorRef{}, err
	}
	return OrchestratorRef{
		env:        env,
		parameters: oParams,
		envVars:    envVars,
		location:   location,
	}, nil
}

func (r *OrchestratorRef) merge(other OrchestratorRef) error {
	var err error
	r.parameters, err = r.parameters.inherit(other.parameters)
	if err != nil {
		return err
	}
	r.envVars, err = r.envVars.inherit(other.envVars)
	if err != nil {
		return err
	}
	return nil
}

//Resolve returns the referenced Orchestrator
func (r OrchestratorRef) Resolve() (Orchestrator, error) {
	orchestrator := r.env.Orchestrator
	params, err := r.parameters.inherit(orchestrator.Parameters)
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
		EnvVars:    envVars,
	}, nil
}
