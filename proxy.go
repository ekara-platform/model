package model

//Proxy represents the proxy definition
type Proxy struct {
	Http    string `yaml:"http_proxy" json:",omitempty"`
	Https   string `yaml:"https_proxy" json:",omitempty"`
	NoProxy string `yaml:"no_proxy" json:",omitempty"`
}

func createProxy(yamlRef yamlProxy) Proxy {
	return Proxy{
		Http:    yamlRef.Http,
		Https:   yamlRef.Https,
		NoProxy: yamlRef.NoProxy,
	}
}

func (r Proxy) inherits(parent Proxy) Proxy {
	res := Proxy{
		Http:    parent.Http,
		Https:   parent.Https,
		NoProxy: parent.NoProxy,
	}
	res.Http = r.Http
	res.Https = r.Https
	res.NoProxy = r.NoProxy
	return res
}
