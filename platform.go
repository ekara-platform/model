package model

import (
	"net/url"
	"os"
	"strings"
)

// Lagoon Platform used to manipulate an environment
type LagoonPlatform struct {
	ComponentBase     *url.URL
	ComponentVersions map[string]Version
	DockerRegistry    *url.URL
	Proxy             Proxy

	Component
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
	lagoonRepository := yamlEnv.Lagoon.Repository
	if lagoonRepository == "" {
		lagoonRepository = LagoonCoreRepository
	}

	base, e := createComponentBase(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "lagoon.componentBase")
	}

	dockerReg, e := createDockerRegistry(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "lagoon.dockerRegistry")
	}

	lagoon := LagoonPlatform{
		ComponentBase:     base,
		ComponentVersions: createComponentMap(vErrs, base, yamlEnv),
		DockerRegistry:    dockerReg,
		Proxy:             createProxy(vErrs, yamlEnv),
	}
	lagoon.Component = createComponent(vErrs, lagoon, "lagoon", lagoonRepository, "")
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

func createProxy(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Proxy {
	httpUrl, e := url.Parse(yamlEnv.Lagoon.Proxy.Http)
	if e != nil {
		vErrs.AddError(e, "proxy.http")
	}
	httpsUrl, e := url.Parse(yamlEnv.Lagoon.Proxy.Https)
	if e != nil {
		vErrs.AddError(e, "proxy.https")
	}
	return Proxy{Http: httpUrl, Https: httpsUrl, NoProxy: yamlEnv.Lagoon.Proxy.NoProxy}
}
