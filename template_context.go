package model

type (
	// TemplateContext the context passed to all ekara templates
	TemplateContext struct {
		// Vars represents accessible descriptor variables,
		Vars Parameters
		// Model represents the environment meta-model (read-only)
		Model TEnvironment
		// Component represents information about the current component
		Component struct {
			// Type of the component
			Type string
			// Name of the component
			Name string
			// Parameters of the component
			Params Parameters
			// Proxy info of the component if any
			Proxy Proxy
			// Environment variables of the component
			EnvVars EnvVars
		}
	}
)

// CreateTemplateContext Returns a template context
func CreateTemplateContext(params Parameters) *TemplateContext {
	return &TemplateContext{
		Vars: params,
	}
}

// CloneTemplateContext deeply clone the given template context
func CloneTemplateContext(other *TemplateContext, cr ComponentReferencer) (*TemplateContext, error) {
	tplC := TemplateContext{
		Vars:  CloneParameters(other.Vars),
		Model: other.Model,
	}
	var target interface{}
	target, e := cr.Component()
	if e != nil {
		return nil, e
	}
	if o, ok := target.(Describable); ok {
		tplC.Component.Type = o.DescType()
		tplC.Component.Name = o.DescName()
	}
	if o, ok := target.(ParametersAware); ok {
		tplC.Component.Params = CloneParameters(o.ParamsInfo())
	}
	if o, ok := target.(ProxyAware); ok {
		tplC.Component.Proxy = o.ProxyInfo()
	}
	if o, ok := target.(EnvVarsAware); ok {
		tplC.Component.EnvVars = o.EnvVarsInfo()
	}
	return &tplC, nil
}

// mergeVars merges others parameters into the template context
func (cc *TemplateContext) mergeVars(others Parameters) {
	cc.Vars = others.inherit(cc.Vars)
}
