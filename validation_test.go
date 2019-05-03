package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an empty yaml file.
//
// The validation must complain about all root elements missing
//- Error: empty environment name @name
//- Error: empty component reference @orchestrator
//- Error: no provider specified @providers
//- Error: no node specified @nodes
//- Warning: no stack specified @stacks
//
// There is no message about a missing ekera platform because it has been defaulted
//
func TestValidationNoContent(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "content", false)
	assert.True(t, vErrs.HasErrors())
	assert.True(t, vErrs.HasWarnings())
	assert.Equal(t, 5, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty environment name", "name"))
	assert.True(t, vErrs.contains(Error, "empty component reference", "orchestrator.component"))
	assert.True(t, vErrs.contains(Error, "no provider specified", "providers"))
	assert.True(t, vErrs.contains(Error, "no node specified", "nodes"))
	assert.True(t, vErrs.contains(Warning, "no stack specified", "stacks"))
}

// Test loading an environment without name.
//
// The validation must complain only about the missing name
//- Error: empty environment name @name
//
func TestValidationNoEnvironmentName(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "environment_name", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty environment name", "name"))
}

// Test loading an environment with an invalid name
//
// The validation must complain only about the invalid name
//- Error: the environment name or the qualifier contains a non alphanumeric character @name|qualifier
//
func TestValidateNoValidName(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/grammar/no_valid_name.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "the environment name or the qualifier contains a non alphanumeric character", "name|qualifier"))
}

// Test loading an environment with an invalid qualifier
//
// The validation must complain only about the invalid qualifier
//- Error: the environment name or the qualifier contains a non alphanumeric character @name|qualifier
//
func TestValidateNoValidQualifier(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/grammar/no_valid_qualifier.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "the environment name or the qualifier contains a non alphanumeric character", "name|qualifier"))
}

// Test loading an environment without nodes.
//
// The validation must complain only about the missing nodes
//- Error: no node specified @nodes
//
func TestValidationNoNodes(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no node specified", "nodes"))
}

// Test loading an environment without providers.
//
// The validation must complain only about the missing providers and the reference
// to a missing provider into the node set specification
//
//- Error: no provider specified @providers
//- Error: reference to unknown provider: aws @nodes.managers.provider
//
func TestValidationNoProviders(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "providers", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no provider specified", "providers"))
	assert.True(t, vErrs.contains(Error, "reference to unknown provider: aws", "nodes.managers.provider"))
}

// Test loading an nodeset referencing an unknown provider.
//
// The validation must complain only about the reference on unknown provider
//
//- Error: reference to unknown provider: dummy @nodes.managers.provider
//
func TestValidationNodesUnknownProvider(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/grammar/nodes_unknown_provider.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "reference to unknown provider: dummy", "nodes.managers.provider"))
}

// Test loading an node set without a reference on a provider.
//
// The validation must complain only about the missing provider reference
//- Error: empty provider reference @nodes.managers.provider
//
func TestValidationNoNodesProvider(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes_provider", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty provider reference", "nodes.managers.provider"))
}

// Test loading an environment without orchestator.
//
// The validation must complain only about the missing orchestator
//- Error: empty component reference @orchestrator
//
func TestValidationNoOrchestrator(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "orchestrator", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty component reference", "orchestrator.component"))
}

// Test loading an environment referencing an unknown orchestrator.
//
// The validation must complain only about the reference on unknown orchestrator
//
//- Error: reference to unknown component: dummy @orchestrator
//
func TestValidationUnknownOrchestrator(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/grammar/unknown_orchestrator.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "reference to unknown component: dummy", "orchestrator.component"))

}

// Test loading an environment without stacks.
//
// The validation must complain only about the missing stacks
//- Warning: no stack specified @stacks
//
func TestValidationNoStacks(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "stacks", true)
	assert.False(t, vErrs.HasErrors())
	assert.True(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Warning, "no stack specified", "stacks"))
}

// Test loading an environment referencing an unknown stack.
//
// The validation must complain only about the reference on unknown stack
//
//- Error: reference to unknown component: dummy @stacks.monitoring.component
//
func TestValidationUnknownStack(t *testing.T) {
	env, e := CreateEnvironment(buildURL(t, "./testdata/yaml/grammar/unknown_stack.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "reference to unknown component: dummy", "stacks.monitoring.component"))

}

// Test loading an task without any playbook .
//
// The validation must complain only about the missing playbook
//
//- Error: empty playbook path @tasks.task1.playbook
//
func TestValidationTasksNoPlayBook(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "task_playbook", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty playbook path", "tasks.task1.playbook"))
}

// Test loading an node set without instances.
//
// The validation must complain only about the instances number being a positive number
//
//- Error: instances must be a positive number @nodes.managers.instances
//
func TestValidationNoNodesInstance(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes_instance", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "instances must be a positive number", "nodes.managers.instances"))
}

// Test loading an node set with no volume names
//
// The validation must complain only about the missing name
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidationNoVolumeName(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "volume_name", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty volume path", "nodes.managers.volumes[0].path"))
}

func testEmptyContent(t *testing.T, name string, onlyWarning bool) (ValidationErrors, Environment) {
	file := fmt.Sprintf("./testdata/yaml/grammar/no_%s.yaml", name)
	env, e := CreateEnvironment(buildURL(t, file), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := validate(t, *env, onlyWarning)
	return vErrs, *env
}

func validate(t *testing.T, env Environment, onlyWarning bool) ValidationErrors {
	vErrs := env.Validate()
	if onlyWarning {
		assert.True(t, vErrs.HasWarnings())
	} else {
		assert.True(t, vErrs.HasErrors())
	}
	return vErrs
}
