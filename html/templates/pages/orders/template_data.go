package orders

import "github.com/aliml92/anor"

type Content struct {
	Orders []anor.Order
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
