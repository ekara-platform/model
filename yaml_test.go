package model

import (
	"testing"

	"log"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestLabelCreate(t *testing.T) {
	assert.Equal(t, []string{"label1", "label2"}, createLabels(&ValidationErrors{}, "label1", "label2").AsStrings())
}

func TestLabelContains(t *testing.T) {
	f := createLabels(&ValidationErrors{}, "label1", "label2")
	assert.Equal(t, true, f.MatchesLabels("label1"))
	assert.Equal(t, true, f.MatchesLabels("label2"))
	assert.Equal(t, false, f.MatchesLabels("label3"))
}

func TestCreateEngineFromHttp(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := parseYamlDescriptor(logger, "https://raw.githubusercontent.com/lagoon-platform/model/master/testdata/yaml/complete_descriptor/lagoon.yaml")
	// no error occurred
	assert.Nil(t, e)
}

func TestCreateEngineFromBadHttp(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	_, e := parseYamlDescriptor(logger, "https://github.com/lagoon-platform/engine/tree/master/testdata/DUMMY.yaml")

	// an error occurred
	assert.NotNil(t, e)

	assert.Equal(t, "HTTP Error getting the environment descriptor , error code 404", e.Error())
}

func TestCreateEngineFromLocal(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, "testdata/yaml/lagoon.yaml")
	assert.Nil(t, e) // no error occurred

	assert.Equal(t, "testEnvironment", yamlEnv.Name)                               // importing file have has precedence
	assert.Equal(t, "This is my awesome Lagoon environment.", yamlEnv.Description) // imported files are merged
	assert.Equal(t, []string{"tag1", "tag2"}, yamlEnv.Labels)
	// FIXME assert.MatchesLabels(t, "task1", "task2", "task3", env.Hooks.Provision.After)        // order matters
}

func TestCreateEngineFromLocalComplexParams(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	yamlEnv, e := parseYamlDescriptor(logger, "testdata/yaml/test/lagoon.yaml")
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
}
