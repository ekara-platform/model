package model

import "errors"

type Provider struct {
	root *Environment
	Parameters
	Component

	Name string
}

type ProviderRef struct {
	Parameters
	provider *Provider
}

func (r ProviderRef) Resolve() Provider {
	// copy provider
	provider := *r.provider

	// Override/complete parameters
	for name, value := range r.Parameters.AsMap() {
		provider.Parameters.add(name, value)
	}

	return provider
}

func createProviders(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Provider {
	res := map[string]Provider{}
	if yamlEnv.Providers == nil || len(yamlEnv.Providers) == 0 {
		vErrs.AddError(errors.New("no provider specified"), "providers")
	} else {
		for name, yamlProvider := range yamlEnv.Providers {
			provider := Provider{
				root:       env,
				Parameters: createParameters(vErrs, yamlProvider.Params),
				Name:       name}

			provider.Component = createComponent(vErrs, env, "providers."+name, yamlProvider.Repository, yamlProvider.Version)

			res[name] = provider
		}
	}
	return res
}

func createProviderRef(vErrs *ValidationErrors, env *Environment, location string, yamlRef yamlRef) ProviderRef {
	if len(yamlRef.Name) == 0 {
		vErrs.AddError(errors.New("empty provider reference"), location)
	} else {
		if val, ok := env.Providers[yamlRef.Name]; ok {
			return ProviderRef{Parameters: createParameters(vErrs, yamlRef.Params), provider: &val}
		} else {
			vErrs.AddError(errors.New("unknown provider reference: "+yamlRef.Name), location+".name")
		}
	}
	return ProviderRef{}
}
