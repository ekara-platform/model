package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsedReference(t *testing.T) {
	refs := CreateUsedReferences()
	refs.add("ref1")
	refs.add("ref2")
	refs.add("ref2")
	// Duplicated ref should not be added
	assert.Len(t, refs.Refs, 2)
	assert.Contains(t, refs.Refs, "ref1", "ref2")

	others := CreateUsedReferences()
	others.add("ref3")
	others.add("ref4")
	others.add("ref5")
	// Duplicated ref already defined into refs should not be added
	others.add("ref2")
	for id := range others.Refs {
		refs.AddReference(id)
	}

	assert.Len(t, refs.Refs, 5)

	assert.True(t, refs.IdUsed("ref1"))
	assert.True(t, refs.IdUsed("ref2"))
	assert.True(t, refs.IdUsed("ref3"))
	assert.True(t, refs.IdUsed("ref4"))
	assert.True(t, refs.IdUsed("ref5"))

	assert.Contains(t, refs.Refs, "ref1", "ref2", "ref3", "ref4", "ref5")

}
