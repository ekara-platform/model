package model

type (
	//Orchestrator specifies the orchestrator used to manage the environment
	Orchestrator struct {
		// The component containing the orchestrator
		cRef componentRef
		// The orchestrator parameters
		Parameters Parameters
		// The Docker parameters
		Docker Parameters
		// The orchestrator environment variables
		EnvVars EnvVars
	}
)

func createOrchestrator(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Orchestrator, error) {
	yamlO := yamlEnv.Orchestrator
	params, err := createParameters(yamlO.Params)
	if err != nil {
		return Orchestrator{}, err
	}
	docker, err := createParameters(yamlO.Docker)
	if err != nil {
		return Orchestrator{}, err
	}
	envVars, err := createEnvVars(yamlO.Env)
	if err != nil {
		return Orchestrator{}, err
	}
	o := Orchestrator{
		cRef:       createComponentRef(env, location.appendPath("component"), yamlO.Component, true),
		Parameters: params,
		Docker:     docker,
		EnvVars:    envVars,
	}
	env.Ekara.tagUsedComponent(o)
	return o, nil
}

func (r Orchestrator) validate() ValidationErrors {
	return ErrorOnInvalid(r.cRef)
}

func (r *Orchestrator) merge(other Orchestrator) error {
	var err error
	err = r.cRef.merge(other.cRef)
	if err != nil {
		return err
	}
	r.Parameters, err = r.Parameters.inherit(other.Parameters)
	if err != nil {
		return err
	}
	r.Docker, err = r.Docker.inherit(other.Docker)
	if err != nil {
		return err
	}
	r.EnvVars, err = r.EnvVars.inherit(other.EnvVars)
	if err != nil {
		return err
	}
	return nil
}

//Component returns the referenced component
func (r Orchestrator) Component() (Component, error) {
	return r.cRef.resolve()
}

//ComponentName returns the referenced component name
func (r Orchestrator) ComponentName() string {
	return r.cRef.ref
}
