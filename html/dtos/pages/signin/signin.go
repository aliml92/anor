package signin

import (
	"github.com/aliml92/anor/html/dtos/partials"
	"html/template"
)

type Base struct {
	partials.Header
	Content
}

type Content struct {
	Message template.HTML
}
