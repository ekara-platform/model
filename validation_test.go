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

type ReContent struct {
	ref refValidationDetails
}

func (r ReContent) validationDetails() refValidationDetails {
	return r.ref
}

func TestMatchingReference(t *testing.T) {
	id := "my_id"
	repo := make(map[string]interface{})
	repo[id] = "blablabla"

	r := refValidationDetails{
		Id:        id,
		Type:      "my_ref",
		Mandatory: true,
		Location:  DescriptorLocation{Path: "my_path"},
		Repo:      repo,
	}

	vErrs := ErrorOnInvalid(ReContent{ref: r})
	assert.False(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
}

func TestUnmatchingReference(t *testing.T) {
	id := "my_id"
	repo := make(map[string]interface{})
	repo[id] = "blablabla"

	r := refValidationDetails{
		Id:        "dummy_id",
		Type:      "my_type",
		Mandatory: true,
		Location:  DescriptorLocation{Path: "my_path"},
		Repo:      repo,
	}

	vErrs := ErrorOnInvalid(ReContent{ref: r})
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "reference to unknown my_type: dummy_id", "my_path"))
}
