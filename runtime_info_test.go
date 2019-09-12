package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRunTimeContext(t *testing.T) {
	c := createRunTimeInfo()
	assert.NotNil(t, c)

	p := Provider{Name: "my_name"}
	c.SetTarget(p)
	assert.Equal(t, c.TargetType, "Provider")
	assert.Equal(t, c.TargetName, "my_name")

}
