package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefixIgnoringCase(t *testing.T) {

	assert.True(t, hasPrefixIgnoringCase("Lorem ipsum dolor sit amet, consectetur adipiscing elit", "lorem"))
	assert.False(t, hasPrefixIgnoringCase("Lorem ipsum dolor sit amet, consectetur adipiscing elit", "Loorem"))
}

func TestHasSuffixIgnoringCase(t *testing.T) {
	assert.True(t, hasSuffixIgnoringCase("Lorem ipsum dolor sit amet, consectetur adipiscing elit", "ELIT"))
	assert.False(t, hasSuffixIgnoringCase("Lorem ipsum dolor sit amet, consectetur adipiscing elit", "Ellit"))
}

func TestUnionStringSlice(t *testing.T) {

	a := []string{"1", "2", "2"}
	b := []string{"2", "3", "4", "4"}

	res := union(a, b)
	assert.NotNil(t, res)

	expected := []string{"1", "2", "3", "4"}
	if assert.Equal(t, len(expected), len(res)) {
		for _, v := range res {
			assert.Contains(t, expected, v)
		}
	}
}
