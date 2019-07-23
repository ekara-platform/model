package model

import (
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
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/grammar/stack_unknown_hook.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := CreatePlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	vErrs := env.Validate()
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

func TestMergeStackHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := StackHook{}
	h.Deploy.Before = append(h.Deploy.Before, task1)
	h.Undeploy.Before = append(h.Undeploy.Before, task1)
	o := StackHook{}
	o.Deploy.Before = append(o.Deploy.Before, task2)
	o.Undeploy.Before = append(o.Undeploy.Before, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
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

func TestMergeStackHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := StackHook{}
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Undeploy.After = append(h.Undeploy.After, task1)
	o := StackHook{}
	o.Deploy.After = append(o.Deploy.After, task2)
	o.Undeploy.After = append(o.Undeploy.After, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
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

func TestMergeStackHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := StackHook{}
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Undeploy.After = append(h.Undeploy.After, task1)

	err := h.merge(h)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Deploy.Before))
	assert.Equal(t, 0, len(h.Undeploy.Before))
	assert.Equal(t, 1, len(h.Deploy.After))
	assert.Equal(t, 1, len(h.Undeploy.After))
	assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
	assert.Equal(t, task1.ref, h.Undeploy.After[0].ref)
}
