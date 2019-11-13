package model

type (
	//Orchestrator specifies the orchestrator used to manage the environment
	Orchestrator struct {
		// The component containing the orchestrator
		cRef componentRef
		// The orchestrator parameters
		Parameters Parameters `yaml:",omitempty"`
		// The orchestrator environment variables
		EnvVars EnvVars `yaml:",omitempty"`
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

func (o Orchestrator) EnvVarsInfo() EnvVars {
	return o.EnvVars
}

func (o Orchestrator) ParamsInfo() Parameters {
	return o.Parameters
}

func (o Orchestrator) validate() ValidationErrors {
	return ErrorOnInvalid(o.cRef)
}

func (o *Orchestrator) customize(with Orchestrator) error {
	var err error
	err = o.cRef.customize(with.cRef)
	if err != nil {
		return err
	}
	o.Parameters = with.Parameters.inherit(o.Parameters)
	o.EnvVars = with.EnvVars.inherit(o.EnvVars)
	return nil
}

//Component returns the referenced component
func (o Orchestrator) Component() (Component, error) {
	return o.cRef.resolve()
}

//ComponentName returns the referenced component name
func (o Orchestrator) ComponentName() string {
	return o.cRef.ref
}

//DescType returns the Describable type of the orchestrator
//  Hardcoded to : "Orchestrator"
func (o Orchestrator) DescType() string {
	return "Orchestrator"
}

//DescName returns the Describable name of the node set
func (o Orchestrator) DescName() string {
	return o.ComponentName()
}
