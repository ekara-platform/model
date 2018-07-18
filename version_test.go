package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncludesVersion(t *testing.T) {
	vErrs := ValidationErrors{}
	assert.True(t, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.True(t, createVersion(&vErrs, "<>", "1.2").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.True(t, createVersion(&vErrs, "<>", "1").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.True(t, createVersion(&vErrs, "<>", "master").IncludesVersion(createVersion(&vErrs, "<>", "master")))
	assert.False(t, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1.2")))
	assert.False(t, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1")))
	assert.False(t, createVersion(&vErrs, "<>", "1.2.4").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.False(t, createVersion(&vErrs, "<>", "1.3.3").IncludesVersion(createVersion(&vErrs, "<>", "1.1.3")))
	assert.False(t, createVersion(&vErrs, "<>", "master").IncludesVersion(createVersion(&vErrs, "<>", "test")))
	assert.False(t, createVersion(&vErrs, "<>", "master").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
}

func TestString(t *testing.T) {
	vErrs := ValidationErrors{}
	assert.Equal(t, "v1.2.3", createVersion(&vErrs, "<>", "1.2.3").String())
	assert.Equal(t, "master", createVersion(&vErrs, "<>", "master").String())
}
