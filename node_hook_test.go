package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an nodeset with unknown hooks
//
// The validation must complain only about 2 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidationNodesUnknownHook(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/grammar/nodes_unknown_hook.yaml"), &TemplateContext{})
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

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "nodes.managers.hooks.provision.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "nodes.managers.hooks.provision.after"))
}

// Test loading an nodeset with valid hooks
func TestValidationNodesKnownHook(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/grammar/nodes_known_hook.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := CreatePlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	vErrs := env.Validate()
	assert.NotNil(t, vErrs)
	assert.False(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 0, len(vErrs.Errors))
}

func TestHasNoTaskNode(t *testing.T) {
	h := NodeHook{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeNodeProvision(t *testing.T) {
	h := NodeHook{}
	h.Provision.Before = append(h.Provision.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterNodeProvision(t *testing.T) {
	h := NodeHook{}
	h.Provision.After = append(h.Provision.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeNodeHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := NodeHook{}
	h.Provision.Before = append(h.Provision.Before, task1)

	o := NodeHook{}
	o.Provision.Before = append(o.Provision.Before, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Provision.Before)) {
		assert.Equal(t, 0, len(h.Provision.After))
		assert.Equal(t, task1.ref, h.Provision.Before[0].ref)
		assert.Equal(t, task2.ref, h.Provision.Before[1].ref)
	}

}

func TestMergeNodeHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := NodeHook{}
	h.Provision.After = append(h.Provision.After, task1)
	o := NodeHook{}
	o.Provision.After = append(o.Provision.After, task2)

	err := h.merge(o)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Provision.After)) {
		assert.Equal(t, 0, len(h.Provision.Before))
		assert.Equal(t, task1.ref, h.Provision.After[0].ref)
		assert.Equal(t, task2.ref, h.Provision.After[1].ref)
	}

}

func TestMergeNodeHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := NodeHook{}
	h.Provision.After = append(h.Provision.After, task1)

	err := h.merge(h)
	assert.Nil(t, err)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Provision.Before))
	assert.Equal(t, 1, len(h.Provision.After))
	assert.Equal(t, task1.ref, h.Provision.After[0].ref)

}
