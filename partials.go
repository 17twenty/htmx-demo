package main

import (
	"fmt"
	"io"
	"text/template"
)

const (
	suffix                = "tpl.html"
	partialTemplateFolder = "partials"
)

// partialEncoder will look for a specific {foo}.tpl.html
// and then execute it by name so only useful for 1:1 mappoings
func partialEncoder(wr io.Writer, partialName string, data any) error {

	// Good to understand how the variables work
	// https://pkg.go.dev/text/template#hdr-Variables
	//
	tmplFile := fmt.Sprintf("%s.%s", partialName, suffix)
	tmpl, err := template.New(tmplFile).ParseFiles(fmt.Sprintf("%s/%s", partialTemplateFolder, tmplFile))
	if err != nil {
		logger.Error("Failed to parse template", "full_path", fmt.Sprintf("%s/%s", partialTemplateFolder, tmplFile), "error", err)
		return err
	}
	err = tmpl.Lookup(partialName).Execute(wr, data)
	if err != nil {
		logger.Error("Failed to execute template", "template", tmpl.Name, "error", err)
		return err
	}
	return nil
}
