package model

import (
	"bytes"
	"text/template"
)

//ApplyTemplate apply the parameters on the template represented by the descriptor content
func ApplyTemplate(u EkUrl, descriptorContent []byte, parameters *TemplateContext) (out bytes.Buffer, err error) {

	// Parse/execute it as a Go template
	out = bytes.Buffer{}
	tpl := template.Must(template.New(u.String()).Option("missingkey=error").Parse(string(descriptorContent)))

	err = tpl.Execute(&out, *parameters)
	if err != nil {
		return
	}

	return
}
