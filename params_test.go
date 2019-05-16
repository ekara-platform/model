package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceConcatenation(t *testing.T) {
	parent := createParameters(map[string]interface{}{
		"slice": []int{1, 2, 3},
	})
	child := createParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, child.inherits(parent)["slice"])
}

func TestSliceMismatch(t *testing.T) {
	parent := createParameters(map[string]interface{}{
		"slice": []string{"a", "b", "c"},
	})
	child := createParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	assert.Equal(t, []int{4, 5, 6}, child.inherits(parent)["slice"])
}

func TestMapMerging(t *testing.T) {
	parent := createParameters(map[string]interface{}{
		"key1": map[string]interface{}{
			"key11": "someValue",
			"key12": "otherValue",
		},
	})
	child := createParameters(map[string]interface{}{
		"key1": map[string]interface{}{
			"key13": "thirdValue",
		},
		"key2": "unrelatedValue",
	})
	res := child.inherits(parent)
	assert.Equal(t, "someValue", (res["key1"]).(map[string]interface{})["key11"])
	assert.Equal(t, "otherValue", (res["key1"]).(map[string]interface{})["key12"])
	assert.Equal(t, "thirdValue", (res["key1"]).(map[string]interface{})["key13"])
	assert.Equal(t, "unrelatedValue", res["key2"])
}

func TestInvalidParameters(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("The code did not panic")
		}
	}()
	createParameters(map[string]interface{}{
		"someKey": map[int]interface{}{
			5: "subValue",
		},
	})
}
