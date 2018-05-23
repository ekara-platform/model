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
	assert.Nil(t, e)
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Parameters)
	assert.Equal(t, 2, len(aws.Parameters.parameters))
	assert.Equal(t, "initial_param1", aws.Parameters.parameters["param1"])
	assert.Equal(t, "initial_param3", aws.Parameters.parameters["param3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	params := managers.Provider.Resolve().Parameters.parameters
	assert.NotNil(t, params)
	assert.Equal(t, 3, len(params))
	assert.Equal(t, "overwritten_param1", params["param1"])
	assert.Equal(t, "new_param2", params["param2"])
	assert.Equal(t, "initial_param3", params["param3"])
}

/*
func TestOverwrittenDocker(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, "testdata/yaml/overwritten/lagoon.yaml")
	assert.Nil(t, e)
	aws := yamlEnv.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Docker)
	assert.Equal(t, 1, len(aws.Docker))
	assert.Equal(t, "initial_docker1", aws.Params["docker1"])

	managers := yamlEnv.Nodes["managers"]
	assert.NotNil(t, managers)
	docker := managers.Docker
	assert.NotNil(t, docker)
	assert.Equal(t, 2, len(docker))
	assert.Equal(t, "overwritten_docker1", docker["docker1"])
	assert.Equal(t, "new_docker2", docker["docker2"])
}
*/
