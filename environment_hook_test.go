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
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/grammar/unknown_global_hook.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := createPlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	assert.NotNil(t, env)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 4, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.create.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.create.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.after"))
}

func TestHasNoTaskEnv(t *testing.T) {
	h := EnvironmentHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeEnvCreate(t *testing.T) {
	h := EnvironmentHooks{}
	h.Create.Before = append(h.Create.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvCreate(t *testing.T) {
	h := EnvironmentHooks{}
	h.Create.After = append(h.Create.After, oneTask)
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

func TestMergeEnvironmentHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Create.Before = append(h.Create.Before, task1)
	h.Deploy.Before = append(h.Deploy.Before, task1)
	o := EnvironmentHooks{}
	o.Create.Before = append(o.Create.Before, task2)
	o.Deploy.Before = append(o.Deploy.Before, task2)

	err := h.customize(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Create.Before)) {
		assert.Equal(t, 0, len(h.Create.After))
		assert.Equal(t, task1.ref, h.Create.Before[0].ref)
		assert.Equal(t, task2.ref, h.Create.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.Before)) {
		assert.Equal(t, 0, len(h.Deploy.After))
		assert.Equal(t, task1.ref, h.Deploy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.Before[1].ref)
	}

}

func TestMergeEnvironmentHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Create.After = append(h.Create.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)
	o := EnvironmentHooks{}
	o.Create.After = append(o.Create.After, task2)
	o.Deploy.After = append(o.Deploy.After, task2)

	err := h.customize(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Create.After)) {
		assert.Equal(t, 0, len(h.Create.Before))
		assert.Equal(t, task1.ref, h.Create.After[0].ref)
		assert.Equal(t, task2.ref, h.Create.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.After)) {
		assert.Equal(t, 0, len(h.Deploy.Before))
		assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.After[1].ref)
	}

}

func TestMergeEnvironmentHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := EnvironmentHooks{}
	h.Create.After = append(h.Create.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)

	err := h.customize(h)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Create.Before))
	assert.Equal(t, 0, len(h.Deploy.Before))
	assert.Equal(t, 1, len(h.Create.After))
	assert.Equal(t, 1, len(h.Deploy.After))
	assert.Equal(t, task1.ref, h.Create.After[0].ref)
	assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
}
