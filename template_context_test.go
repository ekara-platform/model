package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemplateContext(t *testing.T) {
	p, err := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(p))

	pc := CreateContext(p)

	assert.NotNil(t, pc)
	assert.Equal(t, 2, len(pc.Vars))
	va, ok := pc.Vars["key1"]
	assert.True(t, ok)
	assert.Equal(t, va, "value1")

	va, ok = pc.Vars["key2"]
	assert.True(t, ok)
	assert.Equal(t, va, "value2")
}

func TestMergeTemplateContext(t *testing.T) {
	p, _ := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2_owverwritten",
	})

	others, _ := CreateParameters(map[string]interface{}{
		"key2": "value2",
		"key3": "value3",
	})

	pc := CreateContext(p)
	err := pc.MergeVars(others)
	assert.Nil(t, err)
	assert.NotNil(t, pc)
	assert.Equal(t, 3, len(pc.Vars))
	va, ok := pc.Vars["key1"]
	assert.True(t, ok)
	assert.Equal(t, va, "value1")

	// Value should be overwritten
	va, ok = pc.Vars["key2"]
	assert.True(t, ok)
	assert.Equal(t, va, "value2_owverwritten")

	// Missing stuff should be added
	va, ok = pc.Vars["key3"]
	assert.True(t, ok)
	assert.Equal(t, va, "value3")
}
