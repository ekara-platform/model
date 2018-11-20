package model

import (
	_ "encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEmptyContent(t *testing.T, name string, onlyWarning bool) ValidationErrors {
	file := fmt.Sprintf("./testdata/yaml/grammar/no_%s.yaml", name)
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl(file), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, onlyWarning)
	return vErrs
}

func validate(t *testing.T, env Environment, logger *log.Logger, onlyWarning bool) ValidationErrors {
	vErrs := env.Validate()
	if onlyWarning {
		assert.True(t, vErrs.HasWarnings())
	} else {
		assert.True(t, vErrs.HasErrors())
	}
	vErrs.Log(logger)
	return vErrs
}

func TestNoEnvironmentName(t *testing.T) {
	vErrs := testEmptyContent(t, "environment_name", false)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "empty environment name", vErrs.Errors[0].Message)
}

func TestNoNodes(t *testing.T) {
	vErrs := testEmptyContent(t, "nodes", false)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "no node specified", vErrs.Errors[0].Message)
}

func TestNoOrchestrator(t *testing.T) {
	vErrs := testEmptyContent(t, "orchestrator", false)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "empty component reference", vErrs.Errors[0].Message)
	assert.Equal(t, "orchestrator", vErrs.Errors[0].Location.Path)
}

func TestNoStacks(t *testing.T) {
	vErrs := testEmptyContent(t, "stacks", true)
	assert.Equal(t, false, vErrs.HasErrors())
	assert.Equal(t, true, vErrs.HasWarnings())
	assert.Equal(t, Warning, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "no stack specified", vErrs.Errors[0].Message)
}

func TestNoNodesProvider(t *testing.T) {
	vErrs := testEmptyContent(t, "nodes_provider", false)

	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.provider.name", vErrs.Errors[0].Location.Path)
	assert.Equal(t, "empty provider reference", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNoNodesInstance(t *testing.T) {
	vErrs := testEmptyContent(t, "nodes_instance", false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.instances", vErrs.Errors[0].Location.Path)
	assert.Equal(t, "instances must be a positive number", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNoVolumeName(t *testing.T) {
	vErrs := testEmptyContent(t, "volume_name", false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.volumes.path", vErrs.Errors[0].Location.Path)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNodesUnknownProvider(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/nodes_unknown_provider.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.provider.name", vErrs.Errors[0].Location.Path)
	assert.Equal(t, "reference to unknown provider: unknown", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNodesUnknownHook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/nodes_unknown_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))
	testHook(t, "nodes.managers.hooks.provision.before", 0, vErrs)
	testHook(t, "nodes.managers.hooks.provision.after", 1, vErrs)
}

func TestNodesKnownHook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/nodes_known_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	logger.Printf("Validation errors %v", vErrs)
	assert.NotNil(t, vErrs)
	assert.Equal(t, false, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 0, len(vErrs.Errors))
}

func TestTasksNoPlayBook(t *testing.T) {
	vErrs := testEmptyContent(t, "task_playbook", false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "tasks.task1.playbook", vErrs.Errors[0].Location.Path)
	assert.Equal(t, "empty playbook path", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestUnknownGlobalHooks(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/unknown_global_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 10, len(vErrs.Errors))

	testHook(t, "hooks.init.before", 0, vErrs)
	testHook(t, "hooks.init.after", 1, vErrs)
	testHook(t, "hooks.provision.before", 2, vErrs)
	testHook(t, "hooks.provision.after", 3, vErrs)
	testHook(t, "hooks.deploy.before", 4, vErrs)
	testHook(t, "hooks.deploy.after", 5, vErrs)
	testHook(t, "hooks.undeploy.before", 6, vErrs)
	testHook(t, "hooks.undeploy.after", 7, vErrs)
	testHook(t, "hooks.destroy.before", 8, vErrs)
	testHook(t, "hooks.destroy.after", 9, vErrs)
}

func testHook(t *testing.T, msg string, index int, vErrs ValidationErrors) {
	assert.Equal(t, msg, vErrs.Errors[index].Location.Path)
	assert.Equal(t, "reference to unknown task: unknown", vErrs.Errors[index].Message)
	assert.Equal(t, Error, vErrs.Errors[index].ErrorType)
}

func TestNoValidName(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/no_valid_name.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, false)
	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "the environment name or the qualifier contains a non alphanumeric character", vErrs.Errors[0].Message)
}

func TestNoValidQualifier(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/no_valid_qualifier.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, env, logger, false)
	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "the environment name or the qualifier contains a non alphanumeric character", vErrs.Errors[0].Message)
}
