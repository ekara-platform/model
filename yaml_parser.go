package model

import (
	"gopkg.in/yaml.v2"
)

// ParseYamlDescriptor returns an environment based on parsing of the
// descriptor located at the provided URL.
func ParseYamlDescriptor(u EkURL, context *TemplateContext) (env yamlEnvironment, err error) {

	// Read descriptor content
	content, err := u.ReadUrl()
	if err != nil {
		return
	}

	//Parse just the "vars:" section of the descriptor
	tempsVars, err := readEnvironmentVars(content)
	if err != nil {
		return
	}

	//Fill the TemplateContext with the vars content of the descriptor
	err = tempsVars.fillContext(u, context)
	if err != nil {
		return
	}

	// Template the content of the environment descriptor with the freshly
	// parsed vars mixed with the params coming from the launch context.
	out, err := ApplyTemplate(u, content, context)
	if err != nil {
		return
	}

	// Unmarshal the resulting YAML to get an environment
	err = yaml.Unmarshal(out.Bytes(), &env)
	if err != nil {
		return
	}
	return
}

// ParseYamlDescriptorReferences returns an the references, declared and used, into
// the environment based on parsing of the descriptor located at the provided URL.
//
// This parsing must be applied to the main descriptor itself and to all of its parents
func ParseYamlDescriptorReferences(url EkURL, context *TemplateContext) (env EnvironmentReferences, err error) {

	// Read descriptor content
	content, err := url.ReadUrl()
	if err != nil {
		return
	}

	//Parse just the "vars:" section of the descriptor
	tempsVars, err := readEnvironmentVars(content)
	if err != nil {
		return
	}

	//Fill the TemplateContext with the vars content of the descriptor
	err = tempsVars.fillContext(url, context)
	if err != nil {
		return
	}

	// Template the content of the environment descriptor with the freshly
	// parsed vars mixed with the params coming from the launch context.
	out, err := ApplyTemplate(url, content, context)
	if err != nil {
		return
	}

	// Unmarshal the resulting YAML to get only references
	err = yaml.Unmarshal(out.Bytes(), &env)
	if err != nil {
		return
	}
	return
}
