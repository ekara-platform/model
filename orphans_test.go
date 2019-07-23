package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrphans(t *testing.T) {
	refs := CreateOrphans()
	refs.add("ref1", "p")
	refs.add("ref2", "p")
	refs.add("ref2", "p")
	refs.add("ref2", "s")
	// Duplicated ref should not be added
	assert.Len(t, refs.Refs, 3)
	assert.Contains(t, refs.Refs, "ref1-p", "ref2-p", "ref2-s")

	others := CreateOrphans()
	others.add("ref3", "p")
	others.add("ref4", "p")
	others.add("ref5", "p")
	// Duplicated ref already defined into refs should not be added
	others.add("ref2", "p")
	refs.AddAll(*others)

	assert.Len(t, refs.Refs, 6)
	assert.Contains(t, refs.Refs, "ref1-p", "ref2-p", "ref3-p", "ref4-p", "ref5-p", "ref2-s")

	refs.NoMoreAnOrhpan("ref1-p")

	assert.Len(t, refs.Refs, 5)
	assert.Contains(t, refs.Refs, "ref2-p", "ref3-p", "ref4-p", "ref5-p", "ref2-s")

}
