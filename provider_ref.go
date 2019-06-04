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
		templates  Patterns
	}
)

// createProviderRef creates a reference to the provider declared into the yaml reference
func createProviderRef(env *Environment, location DescriptorLocation, yamlRef yamlProviderRef) (providerRef, error) {
	params, err := CreateParameters(yamlRef.Params)
	if err != nil {
		return providerRef{}, err
	}
	envVars, err := createEnvVars(yamlRef.Env)
	if err != nil {
		return providerRef{}, err
	}
	proxy, err := createProxy(yamlRef.Proxy)
	if err != nil {
		return providerRef{}, err
	}
	return providerRef{
		env:        env,
		ref:        yamlRef.Name,
		parameters: params,
		proxy:      proxy,
		envVars:    envVars,
		location:   location,
		mandatory:  true,
		templates:  createPatterns(env, location.appendPath("templates_patterns"), yamlRef.Templates),
	}, nil
}

func (r *providerRef) merge(other providerRef) error {
	var err error
	if r.ref == "" {
		r.ref = other.ref
	}
	r.parameters, err = r.parameters.inherit(other.parameters)
	if err != nil {
		return err
	}
	r.envVars, err = r.envVars.inherit(other.envVars)
	if err != nil {
		return err
	}
	r.proxy, err = r.proxy.inherit(other.proxy)
	if err != nil {
		return err
	}
	r.templates = r.templates.inherit(other.templates)
	return nil
}

func (r providerRef) Resolve() (Provider, error) {
	var err error
	err = ErrorOnInvalid(r)
	if err.(ValidationErrors).HasErrors() {
		return Provider{}, err
	}
	provider := r.env.Providers[r.ref]
	params, err := r.parameters.inherit(provider.Parameters)
	if err != nil {
		return Provider{}, err
	}
	envVars, err := r.envVars.inherit(provider.EnvVars)
	if err != nil {
		return Provider{}, err
	}
	proxy, err := r.proxy.inherit(provider.Proxy)
	if err != nil {
		return Provider{}, err
	}

	return Provider{
		Name:       provider.Name,
		cRef:       provider.cRef,
		Parameters: params,
		EnvVars:    envVars,
		Proxy:      proxy,
		Templates:  r.templates.inherit(provider.Templates)}, nil
}

//reference return a validatable representation of the reference on a provider
func (r providerRef) validationDetails() refValidationDetails {
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

func (r providerRef) Component() (Component, error) {
	p, err := r.Resolve()
	if err != nil {
		return Component{}, err
	}
	return p.cRef.resolve()
}

func (r providerRef) ComponentName() string {
	p, err := r.Resolve()
	if err != nil {
		return ""
	}
	return p.cRef.ref
}
