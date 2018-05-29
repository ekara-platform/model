package model

import (
	"net/url"
	"os"
	"strings"
)

type Settings struct {
	ComponentBase  *url.URL
	DockerRegistry *url.URL
	Proxy          Proxy
}

type Proxy struct {
	Http    *url.URL
	Https   *url.URL
	NoProxy string
}

func createSettings(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Settings {
	base, e := getComponentBase(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "settings.componentBase")
	}
	dockerReg, e := getDockerRegistry(yamlEnv)
	if e != nil {
		vErrs.AddError(e, "settings.dockerRegistry")
	}

	settings := Settings{
		ComponentBase:  base,
		DockerRegistry: dockerReg,
		Proxy:          createProxy(vErrs, yamlEnv)}
	return settings
}

func createProxy(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Proxy {
	httpUrl, e := url.Parse(yamlEnv.Settings.Proxy.Http)
	if e != nil {
		vErrs.AddError(e, "proxy.http")
	}
	httpsUrl, e := url.Parse(yamlEnv.Settings.Proxy.Https)
	if e != nil {
		vErrs.AddError(e, "proxy.https")
	}
	return Proxy{Http: httpUrl, Https: httpsUrl, NoProxy: yamlEnv.Settings.Proxy.NoProxy}
}

func getComponentBase(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultComponentBase

	if yamlEnv.Settings.ComponentBase != "" {
		res = yamlEnv.Settings.ComponentBase
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

func getDockerRegistry(yamlEnv *yamlEnvironment) (*url.URL, error) {
	res := DefaultDockerRegistry
	if yamlEnv.Settings.DockerRegistry != "" {
		res = yamlEnv.Settings.DockerRegistry
	}
	return url.Parse(res)
}
