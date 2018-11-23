package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReContent struct {
	ref validatableReference
}

func (r ReContent) reference() validatableReference {
	return r.ref
}

func TestMatchingReference(t *testing.T) {
	id := "my_id"
	repo := make(map[string]interface{})
	repo[id] = "blablabla"

	r := validatableReference{
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

	r := validatableReference{
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
