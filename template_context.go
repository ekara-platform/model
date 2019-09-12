package model

type (
	//TemplateContext the context passed to all ekara templates
	TemplateContext struct {
		//Vars represents all variable passed into the context,
		// Thoses variables are the ones commit from the CLI
		// parameters file merged with the ones coming from
		// each environment descriptor.
		Vars Parameters
		//Vars represents the Environment definition, in Read Only
		Model TEnvironment
		//Runtime represents the run time details
		RunTimeInfo *RunTimeInfo
	}
)

//CreateContext Returns a template context
func CreateContext(params Parameters) *TemplateContext {
	return &TemplateContext{
		Vars:        params,
		RunTimeInfo: createRunTimeInfo(),
	}
}

//MergeVars merges others parameters into the template context
func (cc *TemplateContext) MergeVars(others Parameters) error {
	var err error
	cc.Vars, err = others.inherit(cc.Vars)
	if err != nil {
		return err
	}

	return nil
}
