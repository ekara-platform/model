package model

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment with unknown global hooks
//
// The validation must complain only about 10 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownGlobalHooks(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/unknown_global_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	//log.Printf("Errors %v: ", vErrs)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 10, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.init.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.init.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.provision.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.provision.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.deploy.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.undeploy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.undeploy.after"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.destroy.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "hooks.destroy.after"))
}
