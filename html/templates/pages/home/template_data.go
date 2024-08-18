package home

import "github.com/aliml92/anor/html/templates/pages/home/components"

type Content struct {
	Featured components.Featured
	Popular  components.Collection
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
