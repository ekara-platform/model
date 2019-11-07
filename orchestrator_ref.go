package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type (
	//OrchestratorRef represents a reference on an Orchestrator
	OrchestratorRef struct {
		parameters Parameters
		envVars    EnvVars
		env        *Environment
		location   DescriptorLocation
	}
)

func (r OrchestratorRef) MarshalYAML() (interface{}, error) {
	b, e :=  yaml.Marshal(&struct {
		Parameters Parameters `yaml:",omitempty"`
		EnvVars    EnvVars `yaml:",omitempty"`
	}{
		Parameters: r.parameters, 
		EnvVars: r.envVars,
	})
	fmt.Printf("--> GBE o returned '%s'", string(b))
	return string(b), e
}

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
