package model

import (
	"net/url"

	"bytes"
	"text/template"
)

//ApplyTemplate apply the parameters on the template represented by the descriptor content
func ApplyTemplate(normalizedUrl *url.URL, descriptorContent []byte, parameters map[string]interface{}) (out bytes.Buffer, err error) {

	// Parse/execute it as a Go template
	out = bytes.Buffer{}
	tpl, err := template.New(normalizedUrl.String()).Parse(string(descriptorContent))
	if err != nil {
		return
	}

	err = tpl.Execute(&out, parameters)
	if err != nil {
		return
	}

	return
}
