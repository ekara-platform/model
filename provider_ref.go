package model

type (

	// providerRef represents a reference to a provider
	providerRef struct {
		ref        string
		parameters Parameters
		envVars    EnvVars
		proxy      Proxy
		env        *Environment
		location   DescriptorLocation
		mandatory  bool
	}
)

//reference return a validatable representation of the reference on a provider
func (r providerRef) reference() validatableReference {
	result := make(map[string]interface{})
	for k, v := range r.env.Providers {
		result[k] = v
	}
	return validatableReference{
		Id:        r.ref,
		Type:      "provider",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

func (r *providerRef) merge(other providerRef) error {
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
	r.proxy = r.proxy.inherits(other.proxy)
	return nil
}

func (r providerRef) Resolve() (Provider, error) {
	validationErrors := ErrorOnInvalid(r)
	if validationErrors.HasErrors() {
		return Provider{}, validationErrors
	}
	provider := r.env.Providers[r.ref]
	return Provider{
		Name:       provider.Name,
		Component:  provider.Component,
		Parameters: r.parameters.inherits(provider.Parameters),
		EnvVars:    r.envVars.inherits(provider.EnvVars),
		Proxy:      r.proxy.inherits(provider.Proxy)}, nil
}

// createProviderRef creates a reference to the provider declared into the yaml reference
func createProviderRef(env *Environment, location DescriptorLocation, yamlRef yamlProviderRef) providerRef {
	return providerRef{
		env:        env,
		ref:        yamlRef.Name,
		parameters: createParameters(yamlRef.Params),
		proxy:      createProxy(yamlRef.Proxy),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
		mandatory:  true,
	}
}

func (r providerRef) ResolveComponent() (Component, error) {
	p, err := r.Resolve()
	if err != nil {
		return Component{}, err
	}
	return p.Component.ResolveComponent()
}
