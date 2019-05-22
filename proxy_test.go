package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProxy(t *testing.T) {
	yaml := yamlProxy{
		Http:    "http_value",
		Https:   "https_value",
		NoProxy: "no_proxy_value",
	}

	p, err := createProxy(yaml)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, "http_value", p.Http)
	assert.Equal(t, "https_value", p.Https)
	assert.Equal(t, "no_proxy_value", p.NoProxy)
}

func TestInheritNoEffect(t *testing.T) {
	p := Proxy{
		Http:    "http_value",
		Https:   "https_value",
		NoProxy: "no_proxy_value",
	}

	other := Proxy{
		Http:    "http_other_value",
		Https:   "https_other_value",
		NoProxy: "no_proxy_other_value",
	}

	p, err := p.inherit(other)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	// Everything was defined in the proxy then inherit has no effect
	assert.Equal(t, "http_value", p.Http)
	assert.Equal(t, "https_value", p.Https)
	assert.Equal(t, "no_proxy_value", p.NoProxy)
}

func TestInherit(t *testing.T) {
	p := Proxy{
		Http:    "",
		Https:   "",
		NoProxy: "",
	}

	other := Proxy{
		Http:    "http_other_value",
		Https:   "https_other_value",
		NoProxy: "no_proxy_other_value",
	}

	p, err := p.inherit(other)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	// Nothing was defined in the proxy then inherit has an effect
	assert.Equal(t, "http_other_value", p.Http)
	assert.Equal(t, "https_other_value", p.Https)
	assert.Equal(t, "no_proxy_other_value", p.NoProxy)
}
