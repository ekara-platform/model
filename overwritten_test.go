package model

import (
	"testing"

	"log"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestOverwrittenParam(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, "./testdata/yaml/overwritten/lagoon.yaml")
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

func TestOverwrittenDocker(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, "./testdata/yaml/overwritten/lagoon.yaml")
	assert.Nil(t, e)
	assert.NotNil(t, env)
	assert.NotNil(t, env.Docker)
	assert.Equal(t, 2, len(env.Docker))
	assert.Equal(t, "initial_docker1", env.Docker["docker1"])
	assert.Equal(t, "initial_docker3", env.Docker["docker3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	dockers := managers.Docker
	assert.NotNil(t, dockers)
	assert.Equal(t, 3, len(dockers))
	assert.Equal(t, "overwritten_docker1", dockers["docker1"])
	assert.Equal(t, "new_docker2", dockers["docker2"])
	assert.Equal(t, "initial_docker3", dockers["docker3"])
}
