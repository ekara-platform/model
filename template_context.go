package model

type (
	//TemplateContext the context passed to all ekara templates
	TemplateContext struct {
		//Vars represents all variable passed into the context,
		// Thoses variables are the ones commit from the CLI
		// parameters file merged with the ones comming from
		// each environment descriptor.
		Vars Parameters
	}
)

//CreateContext Returns a template context
func CreateContext(params Parameters) *TemplateContext {
	return &TemplateContext{
		Vars: params,
	}
}

//Merge others parameters into the template context
func (cc *TemplateContext) MergeVars(others Parameters) error {
	var err error
	cc.Vars, err = cc.Vars.inherit(others)
	if err != nil {
		return err
	}
	return nil
}
