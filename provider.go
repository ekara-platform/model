package model

import "errors"

type Provider struct {
	root       *Environment
	Parameters attributes
	Component

	Name string
}

type ProviderRef struct {
	Parameters attributes `yaml:",inline"`
	provider   *Provider
}

func (p ProviderRef) ProviderName() string {
	return p.provider.Name
}

// ComponentId returns the id of the provider component
func (p ProviderRef) ComponentId() string {
	return p.provider.Component.Id
}

// ComponentId returns the provider component
func (p ProviderRef) Component() Component {
	return p.provider.Component
}

func createProviders(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Provider {
	res := map[string]Provider{}
	if yamlEnv.Providers == nil || len(yamlEnv.Providers) == 0 {
		vErrs.AddError(errors.New("no provider specified"), "providers")
	} else {
		for name, yamlProvider := range yamlEnv.Providers {
			provider := Provider{
				root:       env,
				Parameters: createAttributes(yamlProvider.Params, nil),
				Name:       name,
			}

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
			return ProviderRef{Parameters: createAttributes(yamlRef.Params, val.Parameters), provider: &val}
		} else {
			vErrs.AddError(errors.New("unknown provider reference: "+yamlRef.Name), location+".name")
		}
	}
	return ProviderRef{}
}
