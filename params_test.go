package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceConcatenation(t *testing.T) {
	parent, err := createParameters(map[string]interface{}{
		"slice": []int{1, 2, 3},
	})
	assert.Nil(t, err)
	child, err := createParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	assert.Nil(t, err)
	res, err := child.inherit(parent)
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, res["slice"])
}

func TestSliceMismatch(t *testing.T) {
	parent, err := createParameters(map[string]interface{}{
		"slice": []string{"a", "b", "c"},
	})
	assert.Nil(t, err)
	child, err := createParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	assert.Nil(t, err)
	res, err := child.inherit(parent)
	assert.Nil(t, err)
	assert.Equal(t, []int{4, 5, 6}, res["slice"])
}

func TestMapMerging(t *testing.T) {
	parent, err := createParameters(map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"key11": "someValue",
			"key12": "otherValue",
		},
	})
	assert.Nil(t, err)
	child, err := createParameters(map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"key13": "thirdValue",
		},
		"key2": "unrelatedValue",
	})
	assert.Nil(t, err)
	res, err := child.inherit(parent)
	assert.Nil(t, err)
	assert.Equal(t, "someValue", (res["key1"]).(map[interface{}]interface{})["key11"])
	assert.Equal(t, "otherValue", (res["key1"]).(map[interface{}]interface{})["key12"])
	assert.Equal(t, "thirdValue", (res["key1"]).(map[interface{}]interface{})["key13"])
	assert.Equal(t, "unrelatedValue", res["key2"])
}
