package model

import (
	"net/url"
	"os"
	"path/filepath"
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
	u, e := url.Parse(res)
	if e != nil {
		return nil, e
	}

	// If no scheme was present, assume a file
	if u.Scheme == "" {
		u.Scheme = "file"
	}

	// If file exists locally, resolve its absolute path
	if u.Scheme == "file" {
		if _, e := os.Stat(u.Path); e == nil {
			u.Path, e = filepath.Abs(u.Path)
			if e != nil {
				return nil, e
			}
		}
		if !strings.HasSuffix(u.Path, "/") {
			u.Path = u.Path + "/"
		}
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
