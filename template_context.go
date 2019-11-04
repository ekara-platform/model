package model

type (
	// TemplateContext the context passed to all ekara templates
	TemplateContext struct {
		// Vars represents all variable passed into the context,
		// Thoses variables are the ones commit from the CLI
		// parameters file merged with the ones coming from
		// each environment descriptor.
		Vars Parameters
		// Vars represents the Environment definition, in Read Only
		Model TEnvironment
		//Runtime represents the run time details
		RunTimeInfo *RunTimeInfo
	}
)

// CreateTemplateContext Returns a template context
func CreateTemplateContext(params Parameters) *TemplateContext {
	return &TemplateContext{
		Vars:        params,
		RunTimeInfo: createRunTimeInfo(),
	}
}

// CloneTemplateContext deeply clone the given template context
func CloneTemplateContext(other *TemplateContext, env *Environment) *TemplateContext {
	tplC := TemplateContext{
		Vars:        CloneParameters(other.Vars),
		Model:       CreateTEnvironmentForEnvironment(*env),
		RunTimeInfo: createRunTimeInfo(),
	}
	return &tplC
}

// MergeVars merges others parameters into the template context
func (cc *TemplateContext) MergeVars(others Parameters) {
	cc.Vars = others.inherit(cc.Vars)
}
