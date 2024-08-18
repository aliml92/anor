package components

import "html/template"

type Confirmation struct {
	Message template.HTML
	Email   string
}

func (Confirmation) GetTemplateFilename() string {
	return "signup_confirmation.gohtml"
}
