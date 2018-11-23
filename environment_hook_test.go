package model

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment with unknown global hooks
//
// The validation must complain only about 10 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownGlobalHooks(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/unknown_global_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	//log.Printf("Errors %v: ", vErrs)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 10, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.init.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.init.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.provision.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.provision.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.undeploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.undeploy.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.destroy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.destroy.after"))
}

func TestHasNoTaskEnv(t *testing.T) {
	h := EnvironmentHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeEnvInit(t *testing.T) {
	h := EnvironmentHooks{}
	h.Init.Before = append(h.Init.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvInit(t *testing.T) {
	h := EnvironmentHooks{}
	h.Init.After = append(h.Init.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvProvision(t *testing.T) {
	h := EnvironmentHooks{}
	h.Provision.Before = append(h.Provision.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvProvision(t *testing.T) {
	h := EnvironmentHooks{}
	h.Provision.After = append(h.Provision.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvDeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Deploy.Before = append(h.Deploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvDeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Deploy.After = append(h.Deploy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvUndeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Undeploy.Before = append(h.Undeploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvUndeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Undeploy.After = append(h.Undeploy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvDestroy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Destroy.Before = append(h.Destroy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvDestroy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Destroy.After = append(h.Destroy.After, oneTask)
	assert.True(t, h.HasTasks())
}
