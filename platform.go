package model

import (
	"net/url"
	"os"
	"strings"
)

// Lagoon Platform used to manipulate an environment
type LagoonPlatform struct {
	ComponentBase  *url.URL
	DockerRegistry *url.URL
	Components     map[string]Component
	Component      ComponentRef
}

type Proxy struct {
	Http    *url.URL
	Https   *url.URL
	NoProxy string
}

// createLagoonPlatform create the Lagoon Platform based on the given
// repository and version
//
// The yamlRepoVersion must contains a repository and a version! If the repository
// or the version is missing then a  error will be generated
func createLagoonPlatform(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) LagoonPlatform {
	base, e := createComponentBase(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "lagoon.componentBase")
	}

	dockerReg, e := createDockerRegistry(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "lagoon.dockerRegistry")
	}

	lagoon := LagoonPlatform{
		ComponentBase:  base,
		DockerRegistry: dockerReg,
		Components:     map[string]Component{}}

	// Create all components
	for componentName, yamlComponent := range yamlEnv.Lagoon.Components {
		lagoon.Components[componentName] = createComponent(
			vErrs,
			"lagoon.components."+componentName,
			lagoon.ComponentBase,
			componentName,
			yamlComponent.Repository,
			yamlComponent.Version)
	}

	// Core component defaults if not specified
	var yamlCoreComponent yamlComponent
	var ok bool
	if yamlCoreComponent, ok = yamlEnv.Lagoon.Components[LagoonCoreId]; !ok {
		yamlCoreComponent = yamlComponent{
			Repository: LagoonCoreRepository,
			Version:    ""}
	}

	// (Re-)create core component
	lagoon.Components[LagoonCoreId] = createComponent(
		vErrs,
		"lagoon.components."+LagoonCoreId,
		lagoon.ComponentBase,
		LagoonCoreId,
		yamlCoreComponent.Repository,
		yamlCoreComponent.Version)

	lagoon.Component = createComponentRef(vErrs, lagoon.Components, "lagoon", LagoonCoreId)

	return lagoon
}

func createComponentBase(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultComponentBase

	if yamlEnv.Lagoon.ComponentBase != "" {
		res = yamlEnv.Lagoon.ComponentBase
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
	if yamlEnv.Lagoon.DockerRegistry != "" {
		res = yamlEnv.Lagoon.DockerRegistry
	}
	return url.Parse(res)
}
