package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverwrittenProviderParam(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
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
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
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
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
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

func TestOverwrittenOrchestratorParam(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
	assert.NotNil(t, env)
	assert.NotNil(t, env.Orchestrator)
	assert.NotNil(t, env.Orchestrator.Parameters)
	assert.Equal(t, 2, len(env.Orchestrator.Parameters))
	assert.Equal(t, "param_initial_orchestrator1", env.Orchestrator.Parameters["orchestrator1"])
	assert.Equal(t, "param_initial_orchestrator3", env.Orchestrator.Parameters["orchestrator3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersOrchestrator, e := managers.Orchestrator.Resolve()
	assert.Nil(t, e)
	orchestratorParams := managersOrchestrator.Parameters
	assert.NotNil(t, orchestratorParams)
	assert.Equal(t, 3, len(orchestratorParams))
	assert.Equal(t, "param_overwritten_orchestrator1", orchestratorParams["orchestrator1"])
	assert.Equal(t, "param_new_orchestrator2", orchestratorParams["orchestrator2"])
	assert.Equal(t, "param_initial_orchestrator3", orchestratorParams["orchestrator3"])
}

func TestOverwrittenOrchestratorDocker(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
	assert.NotNil(t, env)
	assert.NotNil(t, env.Orchestrator)
	assert.NotNil(t, env.Orchestrator.Docker)
	assert.Equal(t, 2, len(env.Orchestrator.Docker))
	assert.Equal(t, "docker_initial_orchestrator1", env.Orchestrator.Docker["orchestrator1"])
	assert.Equal(t, "docker_initial_orchestrator3", env.Orchestrator.Docker["orchestrator3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersOrchestrator, e := managers.Orchestrator.Resolve()
	assert.Nil(t, e)
	orchestratorDocker := managersOrchestrator.Docker
	assert.NotNil(t, orchestratorDocker)
	assert.Equal(t, 3, len(orchestratorDocker))
	assert.Equal(t, "docker_overwritten_orchestrator1", orchestratorDocker["orchestrator1"])
	assert.Equal(t, "docker_new_orchestrator2", orchestratorDocker["orchestrator2"])
	assert.Equal(t, "docker_initial_orchestrator3", orchestratorDocker["orchestrator3"])
}

func TestOverwrittenOrchestratorEnv(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/overwritten/ekara.yaml"), MainComponentId, &TemplateContext{})
	assert.Nil(t, e)
	assert.NotNil(t, env)
	assert.NotNil(t, env.Orchestrator)
	assert.NotNil(t, env.Orchestrator.EnvVars)
	assert.Equal(t, 2, len(env.Orchestrator.EnvVars))
	assert.Equal(t, "env_initial_orchestrator1", env.Orchestrator.EnvVars["orchestrator1"])
	assert.Equal(t, "env_initial_orchestrator3", env.Orchestrator.EnvVars["orchestrator3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersOrchestrator, e := managers.Orchestrator.Resolve()
	assert.Nil(t, e)
	orchestratorEnvVars := managersOrchestrator.EnvVars
	assert.NotNil(t, orchestratorEnvVars)
	assert.Equal(t, 3, len(orchestratorEnvVars))
	assert.Equal(t, "env_overwritten_orchestrator1", orchestratorEnvVars["orchestrator1"])
	assert.Equal(t, "env_new_orchestrator2", orchestratorEnvVars["orchestrator2"])
	assert.Equal(t, "env_initial_orchestrator3", orchestratorEnvVars["orchestrator3"])
}

// TODO Add test for TaskRef ans Task and stack
