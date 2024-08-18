package signin

import "html/template"

type Content struct {
	Message template.HTML
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
