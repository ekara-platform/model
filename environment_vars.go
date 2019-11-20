package model

import (
	"gopkg.in/yaml.v2"
)

type (
	// Structure use to deserialize only the environment vars
	yamlEnvironmentVars struct {
		// The descriptor variables
		yamlVars `yaml:",inline"`
	}
)

//Parse just the "vars:" section of the descriptor
func readEnvironmentVars(content []byte) (*yamlEnvironmentVars, error) {
	tempsVars := &yamlEnvironmentVars{}
	err := yaml.Unmarshal(content, tempsVars)
	if err != nil {
		return tempsVars, err
	}
	return tempsVars, nil
}

//Fill the TemplateContext with the vars content of the descriptor.
//The vars content will be templated, using the initial context.
//Once templated it will be merge into the context
func (v yamlEnvironmentVars) fillContext(u EkURL, context *TemplateContext) error {
	varsBytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	out, err := ApplyTemplate(u, varsBytes, context)
	if err != nil {
		return err
	}
	//Parse the templated"vars:" section of the descriptor
	templated, err := readEnvironmentVars(out.Bytes())
	if err != nil {
		return err
	}

	if len(templated.Vars) > 0 {
		context.mergeVars(templated.Vars)
	}
	return nil
}
