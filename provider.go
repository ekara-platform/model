package model

import (
	"encoding/json"
	"errors"
)

type (
	// Provider contains the whole specification of a cloud provider where to
	// create an environemt
	Provider struct {
		// The Name of the provider
		Name string
		// The component containing the provider
		Component componentRef
		// The provider parameters
		Parameters Parameters
		// The provider environment variables
		EnvVars EnvVars
		// The provider proxy
		Proxy Proxy
	}

	//Providers lists all the providers required to build the environemt
	Providers map[string]Provider
)

//DescType returns the Describable type of the provider
//  Hardcoded to : "Provider"
func (r Provider) DescType() string {
	return "Provider"
}

//DescName returns the Describable name of the provider
func (r Provider) DescName() string {
	return r.Name
}

// MarshalJSON returns the serialized content of provider as JSON
func (r Provider) MarshalJSON() ([]byte, error) {
	component, e := r.Component.Resolve()
	if e != nil {
		return nil, e
	}
	return json.Marshal(struct {
		Name       string     `json:",omitempty"`
		Component  string     `json:",omitempty"`
		Parameters Parameters `json:",omitempty"`
		EnvVars    EnvVars    `json:",omitempty"`
		Proxy      Proxy      `json:",omitempty"`
	}{
		Name:       r.Name,
		Component:  component.Id,
		Parameters: r.Parameters,
		EnvVars:    r.EnvVars,
		Proxy:      r.Proxy,
	})
}

func (r Provider) validate() ValidationErrors {
	return ErrorOnInvalid(r.Component)
}

func (r *Provider) merge(other Provider) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated providers (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Component.merge(other.Component); err != nil {
		return err
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
	r.Proxy = r.Proxy.inherits(other.Proxy)
	return nil
}

// createProviders creates all the providers declared into the provided environment
func createProviders(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) Providers {
	res := Providers{}
	for name, yamlProvider := range yamlEnv.Providers {
		providerLocation := location.appendPath(name)
		res[name] = Provider{
			Name:       name,
			Component:  createComponentRef(env, providerLocation.appendPath("component"), yamlProvider.Component, true),
			Parameters: createParameters(yamlProvider.Params),
			Proxy:      createProxy(yamlProvider.Proxy),
			EnvVars:    createEnvVars(yamlProvider.Env)}
	}
	return res
}

func (r Providers) merge(env *Environment, other Providers) error {
	for id, p := range other {
		if provider, ok := r[id]; ok {
			if err := provider.merge(p); err != nil {
				return err
			}
		} else {
			p.Component.env = env
			r[id] = p
		}
	}
	return nil
}
