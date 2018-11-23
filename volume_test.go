package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeVolumeUnrelated(t *testing.T) {
	v := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	v.Parameters["p1"] = "val1"

	o := Volume{
		Path:       "Dummy",
		Parameters: Parameters{},
	}
	o.Parameters["p1"] = "val1"

	err := v.merge(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot merge unrelated volumes (Path != Dummy)")
	}
	assert.Equal(t, 1, len(v.Parameters))
	assert.Contains(t, v.Parameters, "p1")
	assert.Equal(t, v.Parameters["p1"], "val1")
}

func TestMergeVolumeItself(t *testing.T) {
	v := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	v.Parameters["p1"] = "val1"

	err := v.merge(v)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(v.Parameters))
	assert.Contains(t, v.Parameters, "p1")
	assert.Equal(t, v.Parameters["p1"], "val1")
}

func TestMergeVolumeNoUpdate(t *testing.T) {
	v := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	v.Parameters["p1"] = "val1"

	o := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	o.Parameters["p1"] = "val1_updated"

	err := v.merge(o)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(v.Parameters))
	assert.Contains(t, v.Parameters, "p1")
	assert.Equal(t, v.Parameters["p1"], "val1")
}

func TestMergeVolumeAddition(t *testing.T) {
	v := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	v.Parameters["p1"] = "val1"

	o := Volume{
		Path:       "Path",
		Parameters: Parameters{},
	}
	o.Parameters["p1"] = "val1"
	o.Parameters["p2"] = "val2"

	err := v.merge(o)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(v.Parameters))
	assert.Contains(t, v.Parameters, "p1")
	assert.Contains(t, v.Parameters, "p2")
	assert.Equal(t, v.Parameters["p1"], "val1")
	assert.Equal(t, v.Parameters["p2"], "val2")
}
