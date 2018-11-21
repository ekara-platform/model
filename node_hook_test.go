package model

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an nodeset with unknown hooks
//
// The validation must complain only about 2 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//

func TestValidationNodesUnknownHook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/nodes_unknown_hook.yaml"), map[string]interface{}{})

	assert.Nil(t, e)
	vErrs := env.Validate()
	//log.Printf("Errors %v: ", vErrs)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "nodes.managers.hooks.provision.before"))
	assert.True(t, vErrs.contains(Error, "reference to unknown task: unknown", "nodes.managers.hooks.provision.after"))
}

// Test loading an nodeset with valid hooks
func TestValidationNodesKnownHook(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := CreateEnvironment(logger, buildUrl("./testdata/yaml/grammar/nodes_known_hook.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.NotNil(t, vErrs)
	assert.False(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 0, len(vErrs.Errors))
}
