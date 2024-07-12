package signup

import "html/template"

type Confirmation struct {
	Message template.HTML
	Email   string
}
