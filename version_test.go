package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludesVersion(t *testing.T) {
	assert.True(t, createTestVersion("1.2.3").IncludesVersion(createTestVersion("1.2.3")))
	assert.True(t, createTestVersion("1.2").IncludesVersion(createTestVersion("1.2.3")))
	assert.True(t, createTestVersion("1").IncludesVersion(createTestVersion("1.2.3")))
	assert.True(t, createTestVersion("master").IncludesVersion(createTestVersion("master")))
	assert.False(t, createTestVersion("1.2.3").IncludesVersion(createTestVersion("1.2")))
	assert.False(t, createTestVersion("1.2.3").IncludesVersion(createTestVersion("1")))
	assert.False(t, createTestVersion("1.2.4").IncludesVersion(createTestVersion("1.2.3")))
	assert.False(t, createTestVersion("1.3.3").IncludesVersion(createTestVersion("1.1.3")))
	assert.False(t, createTestVersion("master").IncludesVersion(createTestVersion("test")))
	assert.False(t, createTestVersion("master").IncludesVersion(createTestVersion("1.2.3")))
}

func TestString(t *testing.T) {
	assert.Equal(t, "1.2.3", createTestVersion("1.2.3").String())
	assert.Equal(t, "master", createTestVersion("master").String())
}

func TestQualifier(t *testing.T) {
	assert.Equal(t, "1.2.3-beta1", createTestVersion("1.2.3-beta1").String())
}

func createTestVersion(full string) Version {
	version, e := createVersion(full)
	if e != nil {
		panic(e)
	}
	return version
}
