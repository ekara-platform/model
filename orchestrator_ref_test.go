package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getOrchestratorRefOrigin() *OrchestratorRef {
	tr1 := &OrchestratorRef{}
	tr1.envVars = make(map[string]string)
	tr1.envVars["key1"] = "val1_target"
	tr1.envVars["key2"] = "val2_target"
	tr1.parameters = make(map[string]interface{})
	tr1.parameters["key1"] = "val1_target"
	tr1.parameters["key2"] = "val2_target"

	return tr1
}

func getOrchestratorRefOther() OrchestratorRef {
	other := OrchestratorRef{}
	other.envVars = make(map[string]string)
	other.envVars["key2"] = "val2_other"
	other.envVars["key3"] = "val3_other"
	other.parameters = make(map[string]interface{})
	other.parameters["key2"] = "val2_other"
	other.parameters["key3"] = "val3_other"

	return other
}

func checkOrchestratorRefMerge(t *testing.T, ta *OrchestratorRef) {

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

func TestOrchestratorRefMerge(t *testing.T) {
	o := getOrchestratorRefOrigin()
	err := o.customize(getOrchestratorRefOther())
	if assert.Nil(t, err) {
		checkOrchestratorRefMerge(t, o)
	}
}

func TestMergeOrchestratorRefItself(t *testing.T) {
	o := getOrchestratorRefOrigin()
	oi := o
	err := o.customize(*o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}
