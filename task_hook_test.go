package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading a task with unknown hooks
//
// The validation must complain only about 2 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownTaskHooks(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/grammar/task_unknown_hook.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := CreatePlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "tasks.task1.hooks.execute.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "tasks.task1.hooks.execute.after"))
}

func TestHasNoTaskTask(t *testing.T) {
	h := TaskHook{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeTaskExecute(t *testing.T) {
	h := TaskHook{}
	h.Execute.Before = append(h.Execute.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterTaskExecute(t *testing.T) {
	h := TaskHook{}
	h.Execute.After = append(h.Execute.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeTaskHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHook{}
	h.Execute.Before = append(h.Execute.Before, task1)
	o := TaskHook{}
	o.Execute.Before = append(o.Execute.Before, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Execute.Before)) {
		assert.Equal(t, 0, len(h.Execute.After))
		assert.Equal(t, task1.ref, h.Execute.Before[0].ref)
		assert.Equal(t, task2.ref, h.Execute.Before[1].ref)
	}
}

func TestMergeTaskHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHook{}
	h.Execute.After = append(h.Execute.After, task1)
	o := TaskHook{}
	o.Execute.After = append(o.Execute.After, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Execute.After)) {
		assert.Equal(t, 0, len(h.Execute.Before))
		assert.Equal(t, task1.ref, h.Execute.After[0].ref)
		assert.Equal(t, task2.ref, h.Execute.After[1].ref)
	}
}

func TestMergeTaskHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := TaskHook{}
	h.Execute.After = append(h.Execute.After, task1)

	err := h.merge(h)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Execute.Before))
	assert.Equal(t, 1, len(h.Execute.After))
	assert.Equal(t, task1.ref, h.Execute.After[0].ref)
}
