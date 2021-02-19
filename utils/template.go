package utils

import (
	"bytes"
	"html/template"
)

func RenderTemplate(tmpl_file string, data interface{}) (string, error) {
	tmpl, err := template.New("").ParseFiles(tmpl_file)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
