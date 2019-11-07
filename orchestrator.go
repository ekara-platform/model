package model

type (
	//Orchestrator specifies the orchestrator used to manage the environment
	Orchestrator struct {
		// The component containing the orchestrator
		cRef componentRef
		// The orchestrator parameters
		Parameters Parameters `yaml:",omitempty"`
		// The orchestrator environment variables
		EnvVars EnvVars  `yaml:",omitempty"`
	}
)

func createOrchestrator(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Orchestrator, error) {
	yamlO := yamlEnv.Orchestrator
	params, err := CreateParameters(yamlO.Params)
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
		EnvVars:    envVars,
	}

	//env.Ekara.tagUsedComponent(o)
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
	r.Parameters, err = with.Parameters.inherit(r.Parameters)
	if err != nil {
		return err
	}
	r.EnvVars, err = with.EnvVars.inherit(r.EnvVars)
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

//DescType returns the Describable type of the orchestrator
//  Hardcoded to : "Orchestrator"
func (r Orchestrator) DescType() string {
	return "Orchestrator"
}

//DescName returns the Describable name of the node set
func (r Orchestrator) DescName() string {
	return r.ComponentName()
}
