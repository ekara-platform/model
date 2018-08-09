package model

import (
	"errors"
)

type Provider struct {
	// The Name of the provider
	Name string
	// The component containing the provider
	Component ComponentRef
	// The provider parameters
	Parameters Parameters
	// The provider environment variables
	EnvVars EnvVars
}

// Reference to a provider
type ProviderRef struct {
	provider   *Provider
	parameters Parameters
	envVars    EnvVars
}

func (p ProviderRef) Resolve() Provider {
	return Provider{
		Component:  p.provider.Component,
		Parameters: p.parameters.inherit(p.provider.Parameters),
		EnvVars:    p.envVars.inherit(p.provider.EnvVars)}
}

// createProviders creates all the providers declared into the provided environment
func createProviders(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Provider {
	res := map[string]Provider{}
	if yamlEnv.Providers == nil || len(yamlEnv.Providers) == 0 {
		vErrs.AddError(errors.New("no provider specified"), "providers")
	} else {
		for name, yamlProvider := range yamlEnv.Providers {
			res[name] = Provider{
				Name:       name,
				Component:  createComponentRef(vErrs, env.Lagoon.Components, "provider."+name, yamlProvider.Component),
				Parameters: createParameters(yamlProvider.Params),
				EnvVars:    createEnvVars(yamlProvider.Env)}
		}
	}
	return res
}

// createProviderRef creates a reference to the provider declared into the yaml reference
func createProviderRef(vErrs *ValidationErrors, location string, env *Environment, yamlRef yamlProviderRef) ProviderRef {
	if len(yamlRef.Name) == 0 {
		vErrs.AddError(errors.New("empty provider reference"), location)
	} else {
		if val, ok := env.Providers[yamlRef.Name]; ok {
			return ProviderRef{
				provider:   &val,
				parameters: createParameters(yamlRef.Params),
				envVars:    createEnvVars(yamlRef.Env)}
		} else {
			vErrs.AddError(errors.New("unknown provider reference: "+yamlRef.Name), location+".name")
		}
	}
	return ProviderRef{}
}
