package model

import (
	"testing"

	"log"
	"os"

	"github.com/stretchr/testify/assert"
	"strings"
)

// Ignored
func ignoreTestParseFromHttp(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := parseYamlDescriptor(logger, buildUrl("https://raw.githubusercontent.com/lagoon-platform/model/master/testdata/yaml/complete.yaml"), map[string]interface{}{})
	// no error occurred
	assert.Nil(t, e)
}

func TestCreateEngineFromBadHttp(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := parseYamlDescriptor(logger, buildUrl("https://github.com/lagoon-platform/engine/tree/master/testdata/DUMMY.yaml"), map[string]interface{}{})
	// an error occurred
	assert.NotNil(t, e)
	assert.True(t, strings.HasSuffix(e.Error(), "HTTP status 404"))
}

func TestCreateEngineFromLocal(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, buildUrl("testdata/yaml/complete.yaml"), map[string]interface{}{})
	assert.Nil(t, e) // no error occurred

	assert.Equal(t, "testEnvironment", yamlEnv.Name)                               // importing file have has precedence
	assert.Equal(t, "This is my awesome Lagoon environment.", yamlEnv.Description) // imported files are merged
}

func TestCreateEngineFromLocalComplexParams(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, buildUrl("testdata/yaml/complex.yaml"), map[string]interface{}{})
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
}
func TestCreateEngineFromLocalWithData(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, buildUrl("testdata/yaml/data.yaml"), map[string]interface{}{
		"info": map[string]string{
			"name": "Name from data",
			"desc": "Description from data"}})
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
	assert.Equal(t, "Name from data", yamlEnv.Name)
	assert.Equal(t, "Description from data", yamlEnv.Description)
}
