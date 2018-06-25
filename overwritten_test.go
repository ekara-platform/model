package model

import (
	"testing"

	"log"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestOverwrittenParam(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, buildUrl("./testdata/yaml/overwritten/lagoon.yaml"))
	assert.Nil(t, e)
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Parameters)
	assert.Equal(t, 2, len(aws.Parameters))
	assert.Equal(t, "initial_param1", aws.Parameters["param1"])
	assert.Equal(t, "initial_param3", aws.Parameters["param3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	params := managers.Provider.Parameters
	assert.NotNil(t, params)
	assert.Equal(t, 3, len(params))
	assert.Equal(t, "overwritten_param1", params["param1"])
	assert.Equal(t, "new_param2", params["param2"])
	assert.Equal(t, "initial_param3", params["param3"])
}

func TestOverwrittenOrchestratorParam(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, buildUrl("./testdata/yaml/overwritten/lagoon.yaml"))
	assert.Nil(t, e)
	assert.NotNil(t, env)
	assert.NotNil(t, env.Orchestrator)
	assert.NotNil(t, env.Orchestrator.Parameters)
	assert.Equal(t, 2, len(env.Orchestrator.Parameters))
	assert.Equal(t, "param_initial_orchestrator1", env.Orchestrator.Parameters["orchestrator1"])
	assert.Equal(t, "param_initial_orchestrator3", env.Orchestrator.Parameters["orchestrator3"])

	assert.Equal(t, 2, len(env.Orchestrator.Docker))
	assert.Equal(t, "docker_initial_orchestrator1", env.Orchestrator.Docker["orchestrator1"])
	assert.Equal(t, "docker_initial_orchestrator3", env.Orchestrator.Docker["orchestrator3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	orchestrator := managers.Orchestrator.Parameters
	assert.NotNil(t, orchestrator)
	assert.Equal(t, 3, len(orchestrator))
	assert.Equal(t, "param_overwritten_orchestrator1", orchestrator["orchestrator1"])
	assert.Equal(t, "param_new_orchestrator2", orchestrator["orchestrator2"])
	assert.Equal(t, "param_initial_orchestrator3", orchestrator["orchestrator3"])

	orchestrator = managers.Orchestrator.Docker
	assert.NotNil(t, orchestrator)
	assert.Equal(t, 3, len(orchestrator))
	assert.Equal(t, "docker_overwritten_orchestrator1", orchestrator["orchestrator1"])
	assert.Equal(t, "docker_new_orchestrator2", orchestrator["orchestrator2"])
	assert.Equal(t, "docker_initial_orchestrator3", orchestrator["orchestrator3"])
}
