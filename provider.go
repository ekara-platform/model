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

func (r Provider) validate() ValidationErrors {
	return ErrorOnInvalid(r.Component)
}

func (r *Provider) customize(with Provider) error {
	var err error

	if err = r.cRef.customize(with.cRef); err != nil {
		return err
	}

	if r.Name != with.Name {
		return errors.New("cannot customize unrelated providers (" + r.Name + " != " + with.Name + ")")
	}
	if err = r.cRef.customize(with.cRef); err != nil {
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

	r.Proxy, err = r.Proxy.inherit(with.Proxy)
	if err != nil {
		return err
	}
	return nil
}

//Component returns the referenced component
func (r Provider) Component() (Component, error) {
	return r.cRef.resolve()
}

//ComponentName returns the referenced component name
func (r Provider) ComponentName() string {
	return r.cRef.ref
}

// createProviders creates all the providers declared into the provided environment
func createProviders(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (Providers, error) {
	res := Providers{}
	for name, yamlProvider := range yamlEnv.Providers {
		providerLocation := location.appendPath(name)
		params, err := CreateParameters(yamlProvider.Params)
		if err != nil {
			return res, err
		}
		envVars, err := createEnvVars(yamlProvider.Env)
		if err != nil {
			return res, err
		}
		proxy, err := createProxy(yamlProvider.Proxy)
		if err != nil {
			return res, err
		}
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
