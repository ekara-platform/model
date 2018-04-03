package descriptor

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

func createProviders(env *Environment, yamlEnv *yamlEnvironment) (res map[string]Provider, err error) {
	res = map[string]Provider{}
	for name, yamlProvider := range yamlEnv.Providers {
		res[name] = Provider{
			root:       env,
			Parameters: createParameters(yamlProvider.Params),
			Name:       name}
	}
	return
}

func createProviderRef(env *Environment, yamlRef yamlRef) (res ProviderRef, err error) {
	if len(yamlRef.Name) == 0 {
		err = errors.New("empty provider reference")
		return
	}
	if val, ok := env.Providers[yamlRef.Name]; ok {
		res = ProviderRef{Parameters: createParameters(yamlRef.Params), provider: &val}
	} else {
		err = errors.New("unknown provider reference: " + yamlRef.Name)
		return
	}
	return
}
