package model

import (
	"encoding/json"
	"errors"
	"fmt"
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
	// The provider proxy
	Proxy Proxy
}

func (p Provider) HumanDescribe() string {
	return fmt.Sprintf("Provider: %s", p.Name)
}

func (r Provider) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string     `json:",omitempty"`
		Component  string     `json:",omitempty"`
		Parameters Parameters `json:",omitempty"`
		EnvVars    EnvVars    `json:",omitempty"`
		Proxy      Proxy      `json:",omitempty"`
	}{
		Name:       r.Name,
		Component:  r.Component.Resolve().Id,
		Parameters: r.Parameters,
		EnvVars:    r.EnvVars,
		Proxy:      r.Proxy,
	})
}

// Reference to a provider
type ProviderRef struct {
	provider   *Provider
	parameters Parameters
	envVars    EnvVars
	proxy      Proxy
}

func (p ProviderRef) Resolve() Provider {
	return Provider{
		Name:       p.provider.Name,
		Component:  p.provider.Component,
		Parameters: p.parameters.inherit(p.provider.Parameters),
		EnvVars:    p.envVars.inherit(p.provider.EnvVars),
		Proxy:      p.proxy.inherit(p.provider.Proxy),
	}
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
				Component:  createComponentRef(vErrs, env.Ekara.Components, "provider."+name, yamlProvider.Component),
				Parameters: createParameters(yamlProvider.Params),
				Proxy:      createProxy(yamlProvider.Proxy),
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
				proxy:      createProxy(yamlRef.Proxy),
				envVars:    createEnvVars(yamlRef.Env),
			}
		} else {
			vErrs.AddError(errors.New("unknown provider reference: "+yamlRef.Name), location+".name")
		}
	}
	return ProviderRef{}
}
