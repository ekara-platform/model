package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCopies(t *testing.T) {
	copies := []yamlCopy{
		yamlCopy{
			Path:     "path1",
			Patterns: []string{"path1_pattern1", "path1_pattern2"},
		},
		yamlCopy{
			Path:     "path2",
			Patterns: []string{"path2_pattern1", "path2_pattern2"},
		},
	}
	c := createCopies(&Environment{}, DescriptorLocation{Path: "location"}, copies)
	assert.NotNil(t, c)
	assert.Equal(t, len(copies), len(c.Content))

	val, ok := c.Content["path1"]
	if assert.True(t, ok) {
		assert.Contains(t, val.Content, "path1_pattern1")
		assert.Contains(t, val.Content, "path1_pattern2")
	}
	val, ok = c.Content["path2"]
	if assert.True(t, ok) {
		assert.Contains(t, val.Content, "path2_pattern1")
		assert.Contains(t, val.Content, "path2_pattern2")
	}

}

func TestCopiesInherits(t *testing.T) {

	path1 := yamlCopy{
		Path:     "path1",
		Patterns: []string{"path1_pattern1", "path1_pattern2"},
	}
	path2 := yamlCopy{
		Path:     "path2",
		Patterns: []string{"path2_pattern1", "path2_pattern2"},
	}

	c := createCopies(&Environment{}, DescriptorLocation{Path: "location"}, []yamlCopy{path1, path2})

	path1Updated := yamlCopy{
		Path:     "path1",
		Patterns: []string{"path1_pattern1", "path1_patternNew"},
	}

	path3 := yamlCopy{
		Path:     "path3",
		Patterns: []string{"path3_pattern1", "path3_pattern2"},
	}

	o := createCopies(&Environment{}, DescriptorLocation{Path: "location"}, []yamlCopy{path1Updated, path3})

	res := c.inherits(o)
	assert.NotNil(t, res)

	expected := []yamlCopy{
		yamlCopy{
			Path:     "path1",
			Patterns: []string{"path1_pattern1", "path1_pattern2", "path1_patternNew"},
		},
		path2,
		path3,
	}

	if assert.Equal(t, len(expected), len(res.Content)) {

		val, ok := res.Content["path1"]
		if assert.True(t, ok) {
			for _, v := range val.Content {
				assert.Contains(t, expected[0].Patterns, v)
			}
		}

		val, ok = res.Content["path2"]
		if assert.True(t, ok) {
			for _, v := range val.Content {
				assert.Contains(t, expected[1].Patterns, v)
			}
		}

		val, ok = res.Content["path3"]
		if assert.True(t, ok) {
			for _, v := range val.Content {
				assert.Contains(t, expected[2].Patterns, v)
			}
		}
	}
}
