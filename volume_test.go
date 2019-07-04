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
	checkMapInterface(t, v.Parameters, "p1", "val1")
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
	checkMapInterface(t, v.Parameters, "p1", "val1")
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
	checkMapInterface(t, v.Parameters, "p1", "val1")
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
	checkMapInterface(t, v.Parameters, "p1", "val1")
	checkMapInterface(t, v.Parameters, "p2", "val2")
}

func TestMergeNoVolume(t *testing.T) {
	v1 := &Volume{
		Path: "Path1",
	}
	v2 := &Volume{
		Path: "Path2",
	}
	vs := Volumes{}
	vs[v1.Path] = v1
	vs[v2.Path] = v2

	emptyVs := Volumes{}

	vs.merge(emptyVs)
	assert.Equal(t, 2, len(vs))
}

func TestMergeVolumes(t *testing.T) {
	v1 := &Volume{
		Path:       "Path1",
		Parameters: Parameters{},
	}
	v1.Parameters["p1"] = "val1_1"

	v2 := &Volume{
		Path:       "Path2",
		Parameters: Parameters{},
	}
	v2.Parameters["p1"] = "val2_1"

	vs := Volumes{}
	vs[v1.Path] = v1
	vs[v2.Path] = v2

	o1 := &Volume{
		Path:       "Path1",
		Parameters: Parameters{},
	}
	o1.Parameters["p1"] = "update" // Not supposed to be merge
	o1.Parameters["p2"] = "new"    // Must be merged

	o3 := &Volume{ // The whole volume is supposed to be merged
		Path:       "Path3",
		Parameters: Parameters{},
	}
	o3.Parameters["p1"] = "val3_1"
	o3.Parameters["p2"] = "val3_2"

	os := Volumes{}
	os[o1.Path] = o1
	os[o3.Path] = o3

	vs.merge(os)

	assert.Equal(t, 3, len(vs))
	if assert.Equal(t, 2, len(vs[v1.Path].Parameters)) {
		checkMapInterface(t, vs[v1.Path].Parameters, "p1", "val1_1")
		checkMapInterface(t, vs[v1.Path].Parameters, "p2", "new")
	}

	if assert.Equal(t, 1, len(vs[v2.Path].Parameters)) {
		checkMapInterface(t, vs[v2.Path].Parameters, "p1", "val2_1")
	}
	if assert.Equal(t, 2, len(vs[o3.Path].Parameters)) {
		checkMapInterface(t, vs[o3.Path].Parameters, "p1", "val3_1")
		checkMapInterface(t, vs[o3.Path].Parameters, "p2", "val3_2")
	}
}

func TestVolumesAsArray(t *testing.T) {
	v1 := &Volume{
		Path: "Path1",
	}
	v2 := &Volume{
		Path: "Path2",
	}
	vs := Volumes{}
	vs[v1.Path] = v1
	vs[v2.Path] = v2

	arr := vs.AsArray()
	assert.Equal(t, len(arr), len(vs))
}
