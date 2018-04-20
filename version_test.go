package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncludesVersion(t *testing.T) {
	vErrs := ValidationErrors{}
	assert.Equal(t, true, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.Equal(t, true, createVersion(&vErrs, "<>", "1.2").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.Equal(t, true, createVersion(&vErrs, "<>", "1").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.Equal(t, false, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1.2")))
	assert.Equal(t, false, createVersion(&vErrs, "<>", "1.2.3").IncludesVersion(createVersion(&vErrs, "<>", "1")))
	assert.Equal(t, false, createVersion(&vErrs, "<>", "1.2.4").IncludesVersion(createVersion(&vErrs, "<>", "1.2.3")))
	assert.Equal(t, false, createVersion(&vErrs, "<>", "1.3.3").IncludesVersion(createVersion(&vErrs, "<>", "1.1.3")))
}
