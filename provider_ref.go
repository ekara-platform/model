package model

type (

	// Reference to a provider
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

func (r ProviderRef) Reference() Reference {
	result := make(map[string]interface{})
	for k, v := range r.env.Providers {
		result[k] = v
	}
	return Reference{
		Id:        r.ref,
		Type:      "provider",
		Mandatory: r.mandatory,
		Location:  r.location,
		Repo:      result,
	}
}

func (r *ProviderRef) merge(other ProviderRef) error {
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters = r.parameters.inherits(other.parameters)
	r.envVars = r.envVars.inherits(other.envVars)
	r.proxy = r.proxy.inherits(other.proxy)
	return nil
}

func (r ProviderRef) Resolve() (Provider, error) {
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
func createProviderRef(env *Environment, location DescriptorLocation, yamlRef yamlProviderRef) ProviderRef {
	return ProviderRef{
		env:        env,
		ref:        yamlRef.Name,
		parameters: createParameters(yamlRef.Params),
		proxy:      createProxy(yamlRef.Proxy),
		envVars:    createEnvVars(yamlRef.Env),
		location:   location,
		mandatory:  true,
	}
}
