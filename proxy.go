package model

//Proxy represents the proxy definition
type Proxy struct {
	Http    string `yaml:"http_proxy" json:",omitempty"`
	Https   string `yaml:"https_proxy" json:",omitempty"`
	NoProxy string `yaml:"no_proxy" json:",omitempty"`
}

func createProxy(yamlRef yamlProxy) (Proxy, error) {
	return Proxy{
		Http:    yamlRef.Http,
		Https:   yamlRef.Https,
		NoProxy: yamlRef.NoProxy,
	}, nil
}

func (r Proxy) inherit(parent Proxy) (Proxy, error) {
	if r.Http == "" {
		r.Http = parent.Http
	}
	if r.Https == "" {
		r.Https = parent.Https
	}
	if r.NoProxy == "" {
		r.NoProxy = parent.NoProxy
	}
	return r, nil
}
