package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependenciesDependsOnSelf(t *testing.T) {
	deps := []string{"stack_name"}
	d := createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", deps)
	assert.NotNil(t, d)
	// A stack cannot depends on itself
	assert.Equal(t, 0, len(d.Content))
}

func TestDependenciesDependsOnOthers(t *testing.T) {
	deps := []string{"stack_name1", "stack_name2", "stack_name3"}
	d := createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", deps)
	assert.NotNil(t, d)
	if assert.Equal(t, len(deps), len(d.Content)) {
		for _, v := range d.Content {
			assert.Contains(t, deps, v.ref)
		}
	}
}

func TestDependenciesInherits(t *testing.T) {

	deps := []string{"deps1", "deps2"}
	d := createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", deps)

	depsOther := []string{"deps1", "deps3", "deps4"}
	o := createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", depsOther)

	res := d.inherit(o)

	assert.NotNil(t, res)

	expected := []string{"deps1", "deps2", "deps3", "deps4"}
	if assert.Equal(t, len(expected), len(res.Content)) {

		for _, v := range res.Content {
			assert.Contains(t, expected, v.ref)
		}
	}
}
