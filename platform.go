package model

import (
	"net/url"
	"os"
	"strings"
)

// Ekara Platform used to manipulate an environment
type EkaraPlatform struct {
	ComponentBase  *url.URL
	DockerRegistry *url.URL
	Components     map[string]Component
	Component      ComponentRef
}

// createEkaraPlatform create the Ekara Platform based on the given
// repository and version
//
// The yamlRepoVersion must contains a repository and a version! If the repository
// or the version is missing then a  error will be generated
func createEkaraPlatform(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) EkaraPlatform {
	base, e := createComponentBase(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "ekara.componentBase")
	}

	dockerReg, e := createDockerRegistry(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "ekara.dockerRegistry")
	}

	ekara := EkaraPlatform{
		ComponentBase:  base,
		DockerRegistry: dockerReg,
		Components:     map[string]Component{}}

	// Create all components
	for componentName, yamlComponent := range yamlEnv.Ekara.Components {
		component, e := CreateComponent(
			ekara.ComponentBase,
			componentName,
			yamlComponent.Repository,
			yamlComponent.Version)
		if e != nil {
			vErrs.AddError(e, "ekara.components."+componentName)
		} else {
			ekara.Components[componentName] = component
		}
	}

	// Core component defaults if not specified
	var yamlCoreComponent yamlComponent
	var ok bool
	if yamlCoreComponent, ok = yamlEnv.Ekara.Components[EkaraCoreId]; !ok {
		yamlCoreComponent = yamlComponent{
			Repository: EkaraCoreRepository,
			Version:    ""}
	}

	// (Re-)create core component
	coreComponent, e := CreateComponent(
		ekara.ComponentBase,
		EkaraCoreId,
		yamlCoreComponent.Repository,
		yamlCoreComponent.Version)
	if e != nil {
		vErrs.AddError(e, "ekara.components."+EkaraCoreId)
	} else {
		ekara.Components[EkaraCoreId] = coreComponent
	}

	ekara.Component = createComponentRef(vErrs, ekara.Components, "ekara", EkaraCoreId)

	return ekara
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

func createDockerRegistry(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultDockerRegistry
	if yamlEnv.Ekara.DockerRegistry != "" {
		res = yamlEnv.Ekara.DockerRegistry
	}
	return url.Parse(res)
}
