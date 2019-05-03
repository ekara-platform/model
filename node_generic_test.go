package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericNode(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	if assert.Equal(t, len(env.NodeSets), 1) {
		n := env.NodeSets["managers"]
		p, e := n.Provider.Resolve()
		assert.Nil(t, e)
		if val, ok := p.Parameters["generic_param1"]; ok {
			assert.Equal(t, val, "new_generic_param1")
		} else {
			assert.Fail(t, "missing generic param")
		}

		if val, ok := p.EnvVars["generic_env1"]; ok {
			assert.Equal(t, val, "new_generic_env1")
		} else {
			assert.Fail(t, "missing generic env var")
		}

		assert.Equal(t, p.Proxy.NoProxy, "overwritten_aws_no_proxy")
		assert.Equal(t, p.Proxy.Https, "generic_https_proxy")
		assert.Equal(t, p.Proxy.Http, "aws_http_proxy")
	}
}
