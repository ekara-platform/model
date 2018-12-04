package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeStackUnrelated(t *testing.T) {
	sta := Stack{
		Name: "Name",
	}

	o := Stack{
		Name: "Dummy",
	}

	err := sta.merge(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot merge unrelated stacks (Name != Dummy)")
	}
}
