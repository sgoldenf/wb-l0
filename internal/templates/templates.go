package templates

import "text/template"

const templatePath = "./resources/html/index.html"

type TemplateData struct {
	OrderData string
}

func NewTemplateCache() (*template.Template, error) {
	ts, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
