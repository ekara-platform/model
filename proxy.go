package model

import "net/url"

type Proxy struct {
	Http    string
	Https   string
	NoProxy string
}

func createProxy(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) Proxy {
	httpUrl, e := url.Parse(yamlEnv.Proxy.Http)
	if e != nil {
		vErrs.AddError(e, "proxy.http")
	}
	httpsUrl, e := url.Parse(yamlEnv.Proxy.Https)
	if e != nil {
		vErrs.AddError(e, "proxy.http")
	}
	return Proxy{Http: httpUrl.String(), Https: httpsUrl.String(), NoProxy: yamlEnv.Proxy.NoProxy}
}
