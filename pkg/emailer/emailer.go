package emailer

import "context"

type Message struct {
	Subject     string
	To          string
	HtmlContent string
}

type Emailer interface {
	NewMessage(subject string, to string, tmplName string, data interface{}) (Message, error)
	Send(ctx context.Context, message Message) error
}
