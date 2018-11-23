package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	assert.True(t, IsAValidQualifier("aa_bb"))
	assert.True(t, IsAValidQualifier("AA_BB"))
	assert.True(t, IsAValidQualifier("11_22"))
	assert.True(t, IsAValidQualifier("aaAA11_aaBB11"))

	assert.True(t, IsAValidQualifier("aa_BB"))
	assert.True(t, IsAValidQualifier("aa_11"))
	assert.True(t, IsAValidQualifier("11_bb"))
	assert.True(t, IsAValidQualifier("11_BB"))

	assert.True(t, IsAValidQualifier("aaAA11_aaBB11"))

	assert.False(t, IsAValidQualifier("a-b"))
}

// Tests on Envionment
func TestFullQualifiedName(t *testing.T) {
	env := Environment{
		Name:      "ABC",
		Qualifier: "DEF",
	}
	assert.Equal(t, "ABC_DEF", env.QualifiedName().String())
}

func TestPartialQualifiedName(t *testing.T) {
	env := Environment{
		Name: "ABC",
	}
	assert.Equal(t, "ABC", env.QualifiedName().String())
}

func TestEmptyQualifiedName(t *testing.T) {
	env := Environment{}
	// check the zero value
	assert.Equal(t, QualifiedName{}, env.QualifiedName())
	assert.Equal(t, "", env.QualifiedName().String())
}

// Tests on Yaml Envionment
func TestFullQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{
		Name:      "ABC",
		Qualifier: "DEF",
	}
	assert.Equal(t, "ABC_DEF", env.QualifiedName().String())
}

func TestPartialQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{
		Name: "ABC",
	}
	assert.Equal(t, "ABC", env.QualifiedName().String())
}

func TestEmptyQualifiedNameYaml(t *testing.T) {
	env := yamlEnvironment{}
	// check the zero value
	assert.Equal(t, QualifiedName{}, env.QualifiedName())
	assert.Equal(t, "", env.QualifiedName().String())
}

func TestValidQualifiedName(t *testing.T) {
	env := Environment{
		Name:      "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Qualifier: "abcdefghijklmnopqrstuvwxyz",
	}
	assert.True(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "0123456789"
	env.Qualifier = "0123456789"
	assert.True(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "à"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "é"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "ù"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "è"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "ç"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "!"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "-"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "&"
	assert.False(t, env.QualifiedName().ValidQualifiedName())

	env.Name = "#"
	assert.False(t, env.QualifiedName().ValidQualifiedName())
}
