package model

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading s stack with unknown hooks
//
// The validation must complain only about 4 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownStackHooks(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/stack_unknown_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	//log.Printf("Errors %v: ", vErrs)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 4, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "stacks.monitoring.hooks.deploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "stacks.monitoring.hooks.deploy.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "stacks.monitoring.hooks.undeploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "stacks.monitoring.hooks.undeploy.after"))

}

func TestHasNoTaskStack(t *testing.T) {
	h := StackHook{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeStackDeploy(t *testing.T) {
	h := StackHook{}
	h.Deploy.Before = append(h.Deploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterStackDeploy(t *testing.T) {
	h := StackHook{}
	h.Deploy.After = append(h.Deploy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeStackUndeploy(t *testing.T) {
	h := StackHook{}
	h.Undeploy.Before = append(h.Undeploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterStackUndeploy(t *testing.T) {
	h := StackHook{}
	h.Undeploy.After = append(h.Undeploy.After, oneTask)
	assert.True(t, h.HasTasks())
}
