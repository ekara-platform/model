package model

import (
	"gopkg.in/yaml.v2"
)

type (
	// ProviderRef represents a reference to a provider
	ProviderRef struct {
		ref        string
		parameters Parameters
		envVars    EnvVars
		proxy      Proxy
		env        *Environment
		location   DescriptorLocation
		mandatory  bool
	}
)

func (r ProviderRef) MarshalYAML() (interface{}, error) {
	b, e := yaml.Marshal(&struct {
		Ref        string
		Parameters Parameters `yaml:",omitempty"`
		EnvVars    EnvVars    `yaml:",omitempty"`
		Proxy      Proxy      `yaml:",omitempty"`
	}{
		Ref:        r.ref,
		Parameters: r.parameters,
		EnvVars:    r.envVars,
		Proxy:      r.proxy,
	})
	return string(b), e
}

// createProviderRef creates a reference to the provider declared into the yaml reference
func createProviderRef(env *Environment, location DescriptorLocation, yamlRef yamlProviderRef) (ProviderRef, error) {
	return ProviderRef{
		env:        env,
		ref:        yamlRef.Name,
		parameters: CreateParameters(yamlRef.Params),
		proxy:      createProxy(yamlRef.Proxy),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
		mandatory:  true,
	}, nil
}

func (r *ProviderRef) customize(with ProviderRef) error {
	if r.ref == "" {
		r.ref = with.ref
	}
	r.parameters = with.parameters.inherit(r.parameters)
	r.envVars = with.envVars.inherit(r.envVars)
	r.proxy = r.proxy.inherit(with.proxy)
	r.mandatory = with.mandatory
	return nil
}

//Resolve returns the referenced Provider
func (r ProviderRef) Resolve() (Provider, error) {
	var err error
	err = ErrorOnInvalid(r)
	if err.(ValidationErrors).HasErrors() {
		return Provider{}, err
	}
	provider := r.env.Providers[r.ref]
	return Provider{
		Name:       provider.Name,
		cRef:       provider.cRef,
		Parameters: r.parameters.inherit(provider.Parameters),
		EnvVars:    r.envVars.inherit(provider.EnvVars),
		Proxy:      r.proxy.inherit(provider.Proxy),
	}, nil
}

//reference return a validatable representation of the reference on a provider
func (r ProviderRef) validationDetails() refValidationDetails {
	result := make(map[string]interface{})
	for k, v := range r.env.Providers {
		result[k] = v
	}
	return refValidationDetails{
		Id:        r.ref,
		Type:      "provider",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

//Component returns the referenced component
func (r ProviderRef) Component() (Component, error) {
	p, err := r.Resolve()
	if err != nil {
		return Component{}, err
	}
	return p.cRef.resolve()
}

//ComponentName returns the referenced component name
func (r ProviderRef) ComponentName() string {
	p, err := r.Resolve()
	if err != nil {
		return ""
	}
	return p.cRef.ref
}
