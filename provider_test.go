package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderDescType(t *testing.T) {
	p := Provider{}
	assert.Equal(t, p.DescType(), "Provider")
}

func TestProviderDescName(t *testing.T) {
	p := Provider{Name: "my_name"}
	assert.Equal(t, p.DescName(), "my_name")
}
