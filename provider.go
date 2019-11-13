package model

import (
	"errors"
)

type (
	// Provider contains the whole specification of a cloud provider where to
	// create an environemt
	Provider struct {
		// The component containing the provider
		cRef componentRef
		// The Name of the provider
		Name string
		// The provider parameters
		Parameters Parameters `yaml:",omitempty"`
		// The provider environment variables
		EnvVars EnvVars `yaml:",omitempty"`
		// The provider proxy
		Proxy Proxy `yaml:",omitempty"`
	}

	//Providers lists all the providers required to build the environemt
	Providers map[string]Provider
)

func (p Provider) EnvVarsInfo() EnvVars {
	return p.EnvVars
}

func (p Provider) ParamsInfo() Parameters {
	return p.Parameters
}

//ProxyInfo returns the proxy info associated with the provider
func (p Provider) ProxyInfo() Proxy {
	return p.Proxy
}

//DescType returns the Describable type of the provider
func (p Provider) DescType() string {
	return "Provider"
}

//DescName returns the Describable name of the provider
func (p Provider) DescName() string {
	return p.Name
}

func (p Provider) validate() ValidationErrors {
	return ErrorOnInvalid(p.Component)
}

func (p *Provider) customize(with Provider) error {
	var err error

	if err = p.cRef.customize(with.cRef); err != nil {
		return err
	}

	if p.Name != with.Name {
		return errors.New("cannot customize unrelated providers (" + p.Name + " != " + with.Name + ")")
	}
	if err = p.cRef.customize(with.cRef); err != nil {
		return err
	}
	p.Parameters = with.Parameters.inherit(p.Parameters)
	p.EnvVars = with.EnvVars.inherit(p.EnvVars)
	p.Proxy = p.Proxy.inherit(with.Proxy)
	return nil
}

//Component returns the referenced component
func (p Provider) Component() (Component, error) {
	return p.cRef.resolve()
}

//ComponentName returns the referenced component name
func (p Provider) ComponentName() string {
	return p.cRef.ref
}

// createProviders creates all the providers declared into the provided environment
func createProviders(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Providers, error) {
	res := Providers{}
	for name, yamlProvider := range yamlEnv.Providers {
		providerLocation := location.appendPath(name)
		params := CreateParameters(yamlProvider.Params)
		envVars := createEnvVars(yamlProvider.Env)
		proxy := createProxy(yamlProvider.Proxy)
		res[name] = Provider{
			Name:       name,
			cRef:       createComponentRef(env, providerLocation.appendPath("component"), yamlProvider.Component, true),
			Parameters: params,
			EnvVars:    envVars,
			Proxy:      proxy,
		}
		//env.Ekara.tagUsedComponent(res[name])
	}
	return res, nil
}

func (r Providers) customize(env *Environment, with Providers) (Providers, error) {
	res := make(map[string]Provider)
	for k, v := range r {
		res[k] = v
	}
	for id, p := range with {
		if provider, ok := res[id]; ok {
			pm := &provider
			if err := pm.customize(p); err != nil {
				return res, err
			}
			res[id] = *pm
		} else {
			p.cRef.env = env
			res[id] = p
		}
	}
	return res, nil
}
