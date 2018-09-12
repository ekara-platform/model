package model

type Proxy struct {
	Http    string `yaml:"http_proxy" json:",omitempty"`
	Https   string `yaml:"https_proxy" json:",omitempty"`
	NoProxy string `yaml:"no_proxy" json:",omitempty"`
}

func (p Proxy) inherit(parent Proxy) Proxy {
	r := Proxy{
		Http:    parent.Http,
		Https:   parent.Https,
		NoProxy: parent.NoProxy,
	}
	r.Http = p.Http
	r.Https = p.Https
	r.NoProxy = p.NoProxy
	return r
}

// createProxy creates a reference to the provider declared into the yaml reference
func createProxy(yamlRef yamlProxy) Proxy {
	r := Proxy{
		Http:    yamlRef.Http,
		Https:   yamlRef.Https,
		NoProxy: yamlRef.NoProxy,
	}
	return r
}
