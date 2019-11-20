package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverwrittenProviderParam(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := createPlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Parameters)
	assert.Equal(t, 2, len(aws.Parameters))
	assert.Equal(t, "initial_param1", aws.Parameters["param1"])
	assert.Equal(t, "initial_param3", aws.Parameters["param3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve()
	assert.Nil(t, e)
	params := managersProvider.Parameters
	assert.NotNil(t, params)
	assert.Equal(t, 4, len(params))
	assert.Equal(t, "overwritten_param1", params["param1"])
	assert.Equal(t, "new_param2", params["param2"])
	assert.Equal(t, "initial_param3", params["param3"])
}

func TestOverwrittenProviderEnv(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := createPlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.EnvVars)
	assert.Equal(t, 2, len(aws.EnvVars))
	assert.Equal(t, "initial_env1", aws.EnvVars["env1"])
	assert.Equal(t, "initial_env3", aws.EnvVars["env3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve()
	assert.Nil(t, e)
	envs := managersProvider.EnvVars
	assert.NotNil(t, envs)
	assert.Equal(t, 4, len(envs))
	assert.Equal(t, "overwritten_env1", envs["env1"])
	assert.Equal(t, "new_env2", envs["env2"])
	assert.Equal(t, "initial_env3", envs["env3"])
}

func TestOverwrittenProviderProxy(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := createPlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Proxy)
	assert.Equal(t, "", aws.Proxy.Https)
	assert.Equal(t, "aws_http_proxy", aws.Proxy.Http)
	assert.Equal(t, "aws_no_proxy", aws.Proxy.NoProxy)

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve()
	assert.Nil(t, e)
	pr := managersProvider.Proxy
	assert.NotNil(t, pr)
	assert.Equal(t, "aws_http_proxy", pr.Http)
	assert.Equal(t, "generic_https_proxy", pr.Https)
	assert.Equal(t, "overwritten_aws_no_proxy", pr.NoProxy)
}

// TODO Add test for TaskRef ans Task and stack
