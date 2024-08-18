package templates

import (
	"github.com/aliml92/anor/html/templates/shared/header"
)

type TemplateData interface {
	GetTemplateFilename() string
}

type Base struct {
	Header  header.Base
	Content TemplateData
}

type AuthBase struct {
	Content TemplateData
}

type CheckoutBase struct {
	Content TemplateData
}
