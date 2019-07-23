package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReferencedComponentUnique(t *testing.T) {
	refs := CreateReferencedComponents()

	base, err := CreateBase("someBase")
	assert.Nil(t, err)
	rep1, err := CreateRepository(base, "someRep1", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep1,
			Id:         "1",
		},
	})

	rep2, err := CreateRepository(base, "someRep2", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep2,
			Id:         "2",
		},
	})

	// This reference shouldn't be added because another one pointing
	// on the same repository has already been added...
	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep2,
			Id:         "3",
		},
	})

	// Even if the owner is different this reference shouldn't be added because another one pointing
	// on the same repository has already been added...
	refs.AddReference(ReferencedComponent{
		Owner: "owner2",
		Component: Component{
			Repository: rep2,
			Id:         "3",
		},
	})

	// This reference shouldn't be added because another with
	// the same id has already been added...
	rep3, err := CreateRepository(base, "someRep3", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep3,
			Id:         "2",
		},
	})

	// Even if the owner is different this reference shouldn't be added because another with
	// the same id has already been added...
	refs.AddReference(ReferencedComponent{
		Owner: "owner2",
		Component: Component{
			Repository: rep3,
			Id:         "2",
		},
	})

	// Duplicated ref should not be added
	assert.Len(t, refs.Refs, 2)
	assert.Equal(t, refs.Refs[0].Component.Repository.Url.Path(), "someRep1/")
	assert.Equal(t, refs.Refs[1].Component.Repository.Url.Path(), "someRep2/")

	assert.True(t, refs.IdReferenced("1"))
	assert.True(t, refs.IdReferenced("2"))
}

func TestReferencedComponentClean(t *testing.T) {
	refs := CreateReferencedComponents()
	base, err := CreateBase("someBase")
	assert.Nil(t, err)

	rep1, err := CreateRepository(base, "someRep1", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep1,
			Id:         "1",
		},
	})

	rep2, err := CreateRepository(base, "someRep2", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep2,
			Id:         "2",
		},
	})

	rep3, err := CreateRepository(base, "someRep3", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep3,
			Id:         "3",
		},
	})

	rep4, err := CreateRepository(base, "someRep4", "master", "")
	assert.Nil(t, err)

	refs.AddReference(ReferencedComponent{
		Owner: "owner1",
		Component: Component{
			Repository: rep4,
			Id:         "4",
		},
	})

	// Component must be sorted in alphabetical order
	sorted := refs.Sorted()
	assert.Len(t, sorted, 4)
	assert.Equal(t, sorted[0].Id, "1")
	assert.Equal(t, sorted[1].Id, "2")
	assert.Equal(t, sorted[2].Id, "3")
	assert.Equal(t, sorted[3].Id, "4")

	used := CreateUsedReferences()
	used.add("1")
	used.add("3")

	refs.Clean(*used)

	// Only 1 and 3 should be present oncee the clean it's done because
	// 2  and 4 are not within the used components
	assert.Len(t, refs.Refs, 2)
	assert.Equal(t, refs.Refs[0].Component.Repository.Url.Path(), "someRep1/")
	assert.Equal(t, refs.Refs[1].Component.Repository.Url.Path(), "someRep3/")

	assert.True(t, refs.IdReferenced("1"))
	assert.True(t, refs.IdReferenced("3"))
}
