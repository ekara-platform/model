package model

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

type Platform struct {
	Base         *url.URL
	Distribution Component
	Components   map[string]Component
}

func createPlatform(yamlEnv *yamlEnvironment) (Platform, error) {
	components := map[string]Component{}

	// Compute the component base for the environment
	base, e := createComponentBase(yamlEnv)
	if e != nil {
		return Platform{}, errors.New("missing component base: " + e.Error())
	}

	// Create the distribution component (mandatory)
	ekaraRepo := yamlEnv.Ekara.Distribution.Repository
	if ekaraRepo == "" {
		ekaraRepo = EkaraComponentRepo
	}
	ekaraComponent, e := CreateComponent(base, EkaraComponentId, ekaraRepo, yamlEnv.Ekara.Distribution.Version)
	if e != nil {
		return Platform{}, errors.New("invalid distribution: " + e.Error())
	}
	setCredentials(&ekaraComponent, yamlEnv.Ekara.Distribution)

	// Create other components of the environment
	for componentName, yamlComponent := range yamlEnv.Ekara.Components {
		component, e := CreateComponent(
			base,
			componentName,
			yamlComponent.Repository,
			yamlComponent.Version,
			yamlComponent.Imports...)
		if e != nil {
			return Platform{}, errors.New("invalid component " + componentName + ": " + e.Error())
		} else {
			setCredentials(&component, yamlComponent.yamlComponent)
			components[componentName] = component
		}
	}

	return Platform{
		Base:         base,
		Distribution: ekaraComponent,
		Components:   components}, nil
}

func (r Platform) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	for _, c := range r.Components {
		vErrs.merge(ErrorOnInvalid(c))
	}
	return vErrs
}

func (r *Platform) merge(other Platform) error {
	for id, c := range other.Components {
		if _, ok := r.Components[id]; !ok {
			r.Components[id] = c
		}
	}
	return nil
}

func setCredentials(component *Component, yamlComponent yamlComponent) {
	if len(yamlComponent.Auth) > 0 {
		component.Authentication = createParameters(yamlComponent.Auth)
	}
}

func createComponentBase(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultComponentBase

	if yamlEnv != nil && yamlEnv.Ekara.Base != "" {
		res = yamlEnv.Ekara.Base
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
