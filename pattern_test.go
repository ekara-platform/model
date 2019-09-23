package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePatterns(t *testing.T) {
	paths := []string{"my_path1", "my_path2", "my_path3"}
	p := createPatterns(&Environment{}, DescriptorLocation{Path: "location"}, paths)
	assert.NotNil(t, p)
	if assert.Equal(t, len(paths), len(p.Content)) {
		for _, v := range p.Content {
			assert.Contains(t, paths, v)
		}
	}
}

func TestPatternsInherits(t *testing.T) {

	paths := []string{"path1", "path2"}
	p := createPatterns(&Environment{}, DescriptorLocation{Path: "location"}, paths)

	pathsOther := []string{"path1", "path3", "path4"}
	o := createPatterns(&Environment{}, DescriptorLocation{Path: "location"}, pathsOther)

	res := p.inherit(o)

	assert.NotNil(t, res)

	expected := []string{"path1", "path2", "path3", "path4"}
	if assert.Equal(t, len(expected), len(res.Content)) {
		for _, v := range res.Content {
			assert.Contains(t, expected, v)
		}
	}
}
