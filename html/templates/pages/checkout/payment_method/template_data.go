package payment_method

import "github.com/aliml92/anor/html/templates/shared"

type Content struct {
	Stepper shared.Stepper
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
