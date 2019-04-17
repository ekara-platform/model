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
