package model

type Proxy struct {
	Http    string `json:",omitempty"`
	Https   string `json:",omitempty"`
	NoProxy string `json:",omitempty"`
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
