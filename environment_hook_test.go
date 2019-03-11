package model

import (
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
	env, e := CreateEnvironment(buildUrl(t, "./testdata/yaml/grammar/unknown_global_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 8, len(vErrs.Errors))

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

func TestMergeEnvironmentHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Provision.Before = append(h.Provision.Before, task1)
	h.Destroy.Before = append(h.Destroy.Before, task1)
	h.Deploy.Before = append(h.Deploy.Before, task1)
	h.Undeploy.Before = append(h.Undeploy.Before, task1)
	o := EnvironmentHooks{}
	o.Provision.Before = append(o.Provision.Before, task2)
	o.Destroy.Before = append(o.Destroy.Before, task2)
	o.Deploy.Before = append(o.Deploy.Before, task2)
	o.Undeploy.Before = append(o.Undeploy.Before, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Provision.Before)) {
		assert.Equal(t, 0, len(h.Provision.After))
		assert.Equal(t, task1.ref, h.Provision.Before[0].ref)
		assert.Equal(t, task2.ref, h.Provision.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Destroy.Before)) {
		assert.Equal(t, 0, len(h.Destroy.After))
		assert.Equal(t, task1.ref, h.Destroy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.Before)) {
		assert.Equal(t, 0, len(h.Deploy.After))
		assert.Equal(t, task1.ref, h.Deploy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Undeploy.Before)) {
		assert.Equal(t, 0, len(h.Undeploy.After))
		assert.Equal(t, task1.ref, h.Undeploy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Undeploy.Before[1].ref)
	}
}

func TestMergeEnvironmentHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Provision.After = append(h.Provision.After, task1)
	h.Destroy.After = append(h.Destroy.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Undeploy.After = append(h.Undeploy.After, task1)
	o := EnvironmentHooks{}
	o.Provision.After = append(o.Provision.After, task2)
	o.Destroy.After = append(o.Destroy.After, task2)
	o.Deploy.After = append(o.Deploy.After, task2)
	o.Undeploy.After = append(o.Undeploy.After, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Provision.After)) {
		assert.Equal(t, 0, len(h.Provision.Before))
		assert.Equal(t, task1.ref, h.Provision.After[0].ref)
		assert.Equal(t, task2.ref, h.Provision.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Destroy.After)) {
		assert.Equal(t, 0, len(h.Destroy.Before))
		assert.Equal(t, task1.ref, h.Destroy.After[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.After)) {
		assert.Equal(t, 0, len(h.Deploy.Before))
		assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Undeploy.After)) {
		assert.Equal(t, 0, len(h.Undeploy.Before))
		assert.Equal(t, task1.ref, h.Undeploy.After[0].ref)
		assert.Equal(t, task2.ref, h.Undeploy.After[1].ref)
	}
}

func TestMergeEnvironmentHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := EnvironmentHooks{}
	h.Provision.After = append(h.Provision.After, task1)
	h.Destroy.After = append(h.Destroy.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Undeploy.After = append(h.Undeploy.After, task1)

	err := h.merge(h)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Provision.Before))
	assert.Equal(t, 0, len(h.Destroy.Before))
	assert.Equal(t, 0, len(h.Deploy.Before))
	assert.Equal(t, 0, len(h.Undeploy.Before))
	assert.Equal(t, 1, len(h.Provision.After))
	assert.Equal(t, 1, len(h.Destroy.After))
	assert.Equal(t, 1, len(h.Deploy.After))
	assert.Equal(t, 1, len(h.Undeploy.After))
	assert.Equal(t, task1.ref, h.Provision.After[0].ref)
	assert.Equal(t, task1.ref, h.Destroy.After[0].ref)
	assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
	assert.Equal(t, task1.ref, h.Undeploy.After[0].ref)
}
