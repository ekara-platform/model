package model

import (
	"net/url"

	"bytes"
	"text/template"
)

func ApplyTemplate(normalizedUrl *url.URL, descriptorContent []byte, parameters map[string]interface{}) (err error, out bytes.Buffer) {

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
