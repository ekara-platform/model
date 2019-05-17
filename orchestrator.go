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

func createOrchestrator(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	o := Orchestrator{
		cRef:       createComponentRef(env, location.appendPath("component"), yamlO.Component, true),
		Parameters: createParameters(yamlO.Params),
		Docker:     createParameters(yamlO.Docker),
		EnvVars:    createEnvVars(yamlO.Env),
	}
	env.Ekara.tagUsedComponent(o)
	return o
}

func (r Orchestrator) validate() ValidationErrors {
	return ErrorOnInvalid(r.cRef)
}

func (r *Orchestrator) merge(other Orchestrator) error {
	if err := r.cRef.merge(other.cRef); err != nil {
		return err
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.Docker = r.Docker.inherits(other.Docker)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
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
