package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

}
