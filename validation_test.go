package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEmptyContent(t *testing.T, name string, onlyWarning bool) ValidationErrors {
	file := fmt.Sprintf("./testdata/yaml/grammar/no_%s.yaml", name)
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl(file))
	vErrs := assertValidationErrors(t, e, logger, onlyWarning)
	return vErrs
}

func assertValidationErrors(t *testing.T, e error, logger *log.Logger, onlyWarning bool) ValidationErrors {
	assert.NotNilf(t, e, "error should not be nil")
	vErrs, ok := e.(ValidationErrors)
	assert.True(t, ok, "error should be of type ValidationError but was: "+e.Error())
	if onlyWarning {
		assert.True(t, vErrs.HasWarnings())
	} else {
		assert.True(t, vErrs.HasErrors())
	}
	vErrs.Log(logger)
	return vErrs
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
	assert.Equal(t, "no orchestrator specified", vErrs.Errors[0].Message)
}

func TestNoStacks(t *testing.T) {
	vErrs := testEmptyContent(t, "stacks", true)
	assert.Equal(t, false, vErrs.HasErrors())
	assert.Equal(t, true, vErrs.HasWarnings())
	assert.Equal(t, Warning, vErrs.Errors[0].ErrorType)
	assert.Equal(t, "no stack specified", vErrs.Errors[0].Message)
}

func TestNoNodesProvider(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/no_nodes_provider.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.provider", vErrs.Errors[0].Location)
	assert.Equal(t, "empty provider reference", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNoNodesInstance(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/no_nodes_instance.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.instances", vErrs.Errors[0].Location)
	assert.Equal(t, "node set instances must be a positive number", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNoVolumeName(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/no_volume_name.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.volumes.path", vErrs.Errors[0].Location)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNodesUnknownProvider(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/nodes_unknown_provider.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "nodes.managers.provider.name", vErrs.Errors[0].Location)
	assert.Equal(t, "unknown provider reference: unknown", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestNodesUnknownHook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/nodes_unknown_hook.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))
	testHook(t, "nodes.managers.hooks.provision.before", 0, vErrs)
	testHook(t, "nodes.managers.hooks.provision.after", 1, vErrs)
}

func TestStacksUnknownOn(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/stacks_unknown_on.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "stacks.monitoring.on", vErrs.Errors[0].Location)
	assert.Equal(t, "unknown node set reference: unknown", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestTasksNoPlayBook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/no_task_playbook.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "tasks.task1.playbook", vErrs.Errors[0].Location)
	assert.Equal(t, "missing playbook", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestTasksUnknownOn(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/tasks_unknown_on.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

	assert.NotNil(t, vErrs)
	assert.Equal(t, true, vErrs.HasErrors())
	assert.Equal(t, false, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.Equal(t, "tasks.task1.on", vErrs.Errors[0].Location)
	assert.Equal(t, "unknown node set reference: unknown", vErrs.Errors[0].Message)
	assert.Equal(t, Error, vErrs.Errors[0].ErrorType)
}

func TestUnknownGlobalHooks(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := Parse(logger, buildUrl("./testdata/yaml/grammar/unknown_global_hook.yaml"))
	vErrs := assertValidationErrors(t, e, logger, false)

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
	assert.Equal(t, msg, vErrs.Errors[index].Location)
	assert.Equal(t, "unknown task reference: unknown", vErrs.Errors[index].Message)
	assert.Equal(t, Error, vErrs.Errors[index].ErrorType)
}

/*
func TestLocal(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, _ := Parse(logger, buildUrl("./testdata/yaml/test.yaml"))
	assert.NotNil(t, env)
	environmentJson, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", string(environmentJson))
}
*/
