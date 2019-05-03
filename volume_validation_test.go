package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an node set with no volume names
//
// The validation must complain only about the missing name
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidationNoVolumeName(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "volume_name", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty volume path", "nodes.managers.volumes[0].path"))
}
