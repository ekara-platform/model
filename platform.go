package model

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

type Platform struct {
	Component  ComponentRef
	Components map[string]Component
}

func createPlatform(env *Environment, yamlEnv *yamlEnvironment) Platform {
	components := map[string]Component{}

	// Compute the component base for the environment
	base, e := createComponentBase(yamlEnv)
	if e != nil {
		env.errors.addError(e, env.location.appendPath("ekara.componentBase"))
	}

	// Create components of the environment
	for componentName, yamlComponent := range yamlEnv.Ekara.Components {
		component, e := CreateComponent(
			base,
			componentName,
			yamlComponent.Repository,
			yamlComponent.Version,
			yamlComponent.Imports...)
		if e != nil {
			env.errors.addError(e, env.location.appendPath("ekara.components."+componentName))
		} else {
			components[componentName] = component
		}
	}

	// Create core component with default values if none already defined
	if _, ok := components[CoreComponentId]; !ok {
		components[CoreComponentId], e = CreateComponent(
			base,
			CoreComponentId,
			CoreComponentRepo,
			"")
		if e != nil {
			env.errors.addError(errors.New("unable to create core component: "+e.Error()), env.location.appendPath("ekara"))
		}
	}

	return Platform{
		Component:  createComponentRef(env, env.location.appendPath("ekara"), CoreComponentId, true),
		Components: components}
}

func (r Platform) validate() ValidationErrors {
	return r.Component.validate()
}

func (r *Platform) merge(other Platform) error {
	for id, c := range other.Components {
		if _, ok := r.Components[id]; !ok {
			r.Components[id] = c
		}
	}
	return nil
}

func createComponentBase(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultComponentBase

	if yamlEnv != nil && yamlEnv.Ekara.ComponentBase != "" {
		res = yamlEnv.Ekara.ComponentBase
	}

	// If file exists locally, resolve its absolute path and convert it to an URL
	var u *url.URL
	if _, e := os.Stat(res); e == nil {
		u, e = PathToUrl(res)
		if e != nil {
			return nil, e
		}
	} else {
		u, e = url.Parse(res)
		if e != nil {
			return nil, e
		}
	}

	// If no protocol, assume file
	if u.Scheme == "" {
		u.Scheme = "file"
	}

	// Add terminal slash to path if missing
	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	return u, nil
}
