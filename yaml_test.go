package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineFromBadHttp(t *testing.T) {
	_, e := ParseYamlDescriptor(buildURL(t, "https://github.com/ekara-platform/engine/tree/master/testdata/DUMMY.yaml"), &TemplateContext{})
	// an error occurred
	assert.NotNil(t, e)
	assert.True(t, strings.HasSuffix(e.Error(), "HTTP status 404"))
}

func TestCreateEngineFromLocal(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "testdata/yaml/complete.yaml"), &TemplateContext{})
	assert.Nil(t, e) // no error occurred

	assert.Equal(t, "testEnvironment", yamlEnv.Name)                              // importing file have has precedence
	assert.Equal(t, "This is my awesome Ekara environment.", yamlEnv.Description) // imported files are merged
}

func TestCreateEngineFromLocalComplexParams(t *testing.T) {
	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "testdata/yaml/complex.yaml"), &TemplateContext{})
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
}
func TestCreateEngineFromLocalWithData(t *testing.T) {

	vars, _ := CreateParameters(map[string]interface{}{
		"info": map[string]string{
			"name": "Name from data",
			"desc": "Description from data",
		},
	})

	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "testdata/yaml/data.yaml"), &TemplateContext{Vars: vars})
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
	assert.Equal(t, "Name from data", yamlEnv.Name)
	assert.Equal(t, "Description from data", yamlEnv.Description)
}
