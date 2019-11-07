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

func (r *OrchestratorRef) customize(with OrchestratorRef) error {
	var err error
	r.parameters, err = with.parameters.inherit(r.parameters)
	if err != nil {
		return err
	}
	r.envVars, err = with.envVars.inherit(r.envVars)
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
