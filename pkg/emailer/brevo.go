package emailer

import (
	"bytes"
	"context"
	"html/template"

	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"

	"github.com/aliml92/anor/config"
)

type brevoEmailer struct {
	client        *brevo.APIClient
	templateCache *template.Template
	config        *config.EmailConfig
}

func NewBrevoEmailer(cfg *config.EmailConfig, t *template.Template) Emailer {
	// create brevo APIClient
	brevoCfg := brevo.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", cfg.APIKey)
	brevoAPIClient := brevo.NewAPIClient(brevoCfg)

	return &brevoEmailer{
		client:        brevoAPIClient,
		templateCache: t,
		config:        cfg,
	}
}

func (e *brevoEmailer) NewMessage(subject string, to string, tmplName string, data interface{}) (Message, error) {
	var buf bytes.Buffer
	err := e.templateCache.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		return Message{}, err
	}
	msg := Message{
		Subject:     subject,
		To:          to,
		HtmlContent: buf.String(),
	}
	return msg, nil
}

func (e *brevoEmailer) Send(ctx context.Context, m Message) error {
	_, _, err := e.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Anor",
			Email: e.config.FromEmail,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: m.To,
			},
		},
		Subject:     m.Subject,
		HtmlContent: m.HtmlContent,
	})

	return err
}
