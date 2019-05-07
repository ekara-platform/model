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
	//                      / |       | \
	//                     4  5       10  11

	sts["6"] = Stack{Name: "6", DependsOn: getDependsOn("2")}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}
	sts["4"] = Stack{Name: "4", DependsOn: getDependsOn("3")}
	sts["9"] = Stack{Name: "9", DependsOn: getDependsOn("8")}
	sts["5"] = Stack{Name: "5", DependsOn: getDependsOn("3")}
	sts["7"] = Stack{Name: "7", DependsOn: getDependsOn("1")}
	sts["8"] = Stack{Name: "8", DependsOn: getDependsOn("1")}
	sts["10"] = Stack{Name: "10", DependsOn: getDependsOn("9")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["12"] = Stack{Name: "12", DependsOn: getDependsOn("8")}
	sts["11"] = Stack{Name: "11", DependsOn: getDependsOn("9")}
	sts["0"] = Stack{Name: "0"}

	assert.Equal(t, 13, len(sts))
	resolved, err := sts.ResolveDependencies()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.DependsOn {
			if _, ok := processed[d.ref]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d.ref)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Equal(t, 13, len(sts))

}

func TestStacksMultiplesDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / | /    | \
	//                       3  6     9  12
	//                      / |   \   | \
	//                     4  5    \-10  11

	sts["6"] = Stack{Name: "6", DependsOn: getDependsOn("2", "7")}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}
	sts["4"] = Stack{Name: "4", DependsOn: getDependsOn("3")}
	sts["9"] = Stack{Name: "9", DependsOn: getDependsOn("8")}
	sts["5"] = Stack{Name: "5", DependsOn: getDependsOn("3")}
	sts["7"] = Stack{Name: "7", DependsOn: getDependsOn("1")}
	sts["8"] = Stack{Name: "8", DependsOn: getDependsOn("1")}
	sts["10"] = Stack{Name: "10", DependsOn: getDependsOn("9", "6")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["12"] = Stack{Name: "12", DependsOn: getDependsOn("8")}
	sts["11"] = Stack{Name: "11", DependsOn: getDependsOn("9")}
	sts["0"] = Stack{Name: "0"}

	assert.Equal(t, 13, len(sts))
	resolved, err := sts.ResolveDependencies()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.DependsOn {
			if _, ok := processed[d.ref]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d.ref)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Equal(t, 13, len(sts))

}

func TestStacksCyclicDependencies(t *testing.T) {
	sts := Stacks{}

	sts["1"] = Stack{Name: "1", DependsOn: getDependsOn("3")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}

	assert.Equal(t, 3, len(sts))
	_, err := sts.ResolveDependencies()
	assert.NotNil(t, err)

	//Check that the original Stacks has been untouched
	assert.Equal(t, 3, len(sts))

}

func getDependsOn(deps ...string) []stackRef {
	res := make([]stackRef, 0)
	for _, s := range deps {
		r := stackRef{
			ref: s,
		}
		res = append(res, r)
	}
	return res
}
