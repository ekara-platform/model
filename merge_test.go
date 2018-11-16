package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeVolume(t *testing.T) {
	paramV1 := make(map[string]interface{})
	paramV1["paramV1_key1"] = "paramV1_val1"
	paramV1["paramV1_key2"] = "paramV1_val2"

	v1 := Volume{
		Name:       "name1",
		Parameters: paramV1,
	}

	paramV2 := make(map[string]interface{})
	paramV2["paramV2_key1"] = "paramV2_val1"
	paramV2["paramV2_key2"] = "paramV2_val2"

	v2 := Volume{
		Name:       "name1",
		Parameters: paramV2,
	}

	v1.merge(v2)
	assert.Equal(t, v1.Name, "name1")
	assert.Equal(t, 4, len(v1.Parameters))

}

/*
func TestMergeVolumeNotRelated(t *testing.T) {
	paramV1 := make(map[string]interface{})
	paramV1["paramV1_key1"] = "paramV1_val1"
	paramV1["paramV1_key2"] = "paramV1_val2"

	v1 := Volume{
		Name:       "name1",
		Parameters: paramV1,
	}

	paramV2 := make(map[string]interface{})
	paramV2["paramV2_key1"] = "paramV2_val1"
	paramV2["paramV2_key2"] = "paramV2_val2"

	v2 := Volume{
		Name:       "name2",
		Parameters: paramV2,
	}

	v1.merge(v2)
	assert.Equal(t, v1.Name, "name1")
	assert.Equal(t, 4, len(v1.Parameters))

}

func TestMergeVolumeVariaousNumber(t *testing.T) {
	paramV1 := make(map[string]interface{})
	paramV1["paramV1_key1"] = "paramV1_val1"
	paramV1["paramV1_key2"] = "paramV1_val2"

	v1 := Volume{
		Name:       "name1",
		Parameters: paramV1,
	}

	paramV2 := make(map[string]interface{})
	paramV2["paramV2_key1"] = "paramV2_val1"
	paramV2["paramV2_key2"] = "paramV2_val2"

	v2 := Volume{
		Name:       "name1",
		Parameters: paramV2,
	}

	paramV3 := make(map[string]interface{})
	paramV3["paramV3_key1"] = "paramV3_val1"
	paramV3["paramV3_key2"] = "paramV3_val2"

	v3 := Volume{
		Name:       "name3",
		Parameters: paramV3,
	}

	orther := make([]Volume, 2)
	orther[0] = v2
	orther[2] = v3

	dest := make([]Volume, 1)
	dest[0] = v1

	// Same loop than the one currently invoked into node.go
	for i, v := range orther {
		dest[i].merge(v)
	}

	v1.merge(v2)
	assert.Equal(t, v1.Name, "name1")
	assert.Equal(t, 4, len(v1.Parameters))

}
*/
