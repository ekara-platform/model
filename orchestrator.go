package model

type (
	//Orchestrator specifies the orchestrator used to manage the environment
	Orchestrator struct {
		// The component containing the orchestrator
		cRef componentRef
		// The orchestrator parameters
		Parameters Parameters
		// The orchestrator environment variables
		EnvVars EnvVars
	}
)

func createOrchestrator(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Orchestrator, error) {
	yamlO := yamlEnv.Orchestrator
	o := Orchestrator{
		cRef:       createComponentRef(env, location.appendPath("component"), yamlO.Component, true),
		Parameters: CreateParameters(yamlO.Params),
		EnvVars:    createEnvVars(yamlO.Env),
	}
	return o, nil
}

func (r Orchestrator) validate() ValidationErrors {
	return ErrorOnInvalid(r.cRef)
}

func (r *Orchestrator) customize(with Orchestrator) error {
	var err error
	err = r.cRef.customize(with.cRef)
	if err != nil {
		return err
	}
	r.Parameters = with.Parameters.inherit(r.Parameters)
	r.EnvVars = with.EnvVars.inherit(r.EnvVars)
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

//DescType returns the Describable type of the orchestrator
//  Hardcoded to : "Orchestrator"
func (r Orchestrator) DescType() string {
	return "Orchestrator"
}

//DescName returns the Describable name of the node set
func (r Orchestrator) DescName() string {
	return r.ComponentName()
}
