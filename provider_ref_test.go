package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getProviderRefOrigin(ref string) *ProviderRef {
	tr1 := &ProviderRef{
		ref:       ref,
		mandatory: true,
	}
	tr1.envVars = make(map[string]string)
	tr1.envVars["key1"] = "val1_target"
	tr1.envVars["key2"] = "val2_target"
	tr1.parameters = make(map[string]interface{})
	tr1.parameters["key1"] = "val1_target"
	tr1.parameters["key2"] = "val2_target"

	return tr1
}

func getProviderRefOther(ref string) *ProviderRef {
	other := &ProviderRef{
		ref:       ref,
		mandatory: false,
	}
	other.envVars = make(map[string]string)
	other.envVars["key2"] = "val2_other"
	other.envVars["key3"] = "val3_other"
	other.parameters = make(map[string]interface{})
	other.parameters["key2"] = "val2_other"
	other.parameters["key3"] = "val3_other"

	return other
}

func checkProviderRefMerge(t *testing.T, ta *ProviderRef, expectedRef string) {

	assert.Equal(t, ta.ref, expectedRef)
	assert.False(t, ta.mandatory)

	if assert.Len(t, ta.envVars, 3) {
		checkMap(t, ta.envVars, "key1", "val1_target")
		checkMap(t, ta.envVars, "key2", "val2_other")
		checkMap(t, ta.envVars, "key3", "val3_other")
	}

	if assert.Len(t, ta.parameters, 3) {
		checkMapInterface(t, ta.parameters, "key1", "val1_target")
		checkMapInterface(t, ta.parameters, "key2", "val2_other")
		checkMapInterface(t, ta.parameters, "key3", "val3_other")
	}
}

func TestProviderRefMerge(t *testing.T) {
	o := getProviderRefOrigin("my_ref")
	err := o.customize(*getProviderRefOther("my_name"))
	if assert.Nil(t, err) {
		checkProviderRefMerge(t, o, "my_ref")
	}
}

func TestProviderRefMergeEmptyLocation(t *testing.T) {
	o := getProviderRefOrigin("")
	err := o.customize(*getProviderRefOther("other_ref"))
	if assert.Nil(t, err) {
		checkProviderRefMerge(t, o, "other_ref")
	}
}

func TestMergeProviderRefItself(t *testing.T) {
	o := getProviderRefOrigin("my_ref")
	oi := o
	err := o.customize(*o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}
