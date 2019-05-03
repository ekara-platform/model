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

func TestStacksDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / |     | \
	//                       3  6     9  12
	//                      / |        | \
	//                     4  5       10  11

	sts["6"] = Stack{Name: "6", DependsOn: "2"}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", DependsOn: "2"}
	sts["4"] = Stack{Name: "4", DependsOn: "3"}
	sts["9"] = Stack{Name: "9", DependsOn: "8"}
	sts["5"] = Stack{Name: "5", DependsOn: "3"}
	sts["7"] = Stack{Name: "7", DependsOn: "1"}
	sts["8"] = Stack{Name: "8", DependsOn: "1"}
	sts["10"] = Stack{Name: "10", DependsOn: "9"}
	sts["2"] = Stack{Name: "2", DependsOn: "1"}
	sts["12"] = Stack{Name: "12", DependsOn: "8"}
	sts["11"] = Stack{Name: "11", DependsOn: "9"}
	sts["0"] = Stack{Name: "0"}

	assert.Equal(t, 13, len(sts))
	ch := sts.ResolveDependencies()

	assert.Equal(t, "0", (<-ch).Name)
	assert.Equal(t, "1", (<-ch).Name)
	assert.Equal(t, "2", (<-ch).Name)
	assert.Equal(t, "3", (<-ch).Name)
	assert.Equal(t, "4", (<-ch).Name)
	assert.Equal(t, "5", (<-ch).Name)
	assert.Equal(t, "6", (<-ch).Name)
	assert.Equal(t, "7", (<-ch).Name)
	assert.Equal(t, "8", (<-ch).Name)
	assert.Equal(t, "9", (<-ch).Name)
	assert.Equal(t, "10", (<-ch).Name)
	assert.Equal(t, "11", (<-ch).Name)
	assert.Equal(t, "12", (<-ch).Name)

	//Check that the original Stacks has been untouched
	assert.Equal(t, 13, len(sts))
}
