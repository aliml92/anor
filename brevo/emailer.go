package brevo

import (
	"bytes"
	"context"
	"html/template"

	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"

	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/email"
)

var _ email.Emailer = (*Emailer)(nil)

type Emailer struct {
	client        *brevo.APIClient
	templateCache *template.Template
	config        *config.EmailConfig
}

func NewEmailer(cfg config.EmailConfig) *Emailer {
	t := template.Must(template.ParseGlob(cfg.Templates))
	brevoCfg := brevo.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", cfg.APIKey)
	brevoClient := brevo.NewAPIClient(brevoCfg)

	return &Emailer{
		client:        brevoClient,
		templateCache: t,
		config:        &cfg,
	}
}

func (e *Emailer) NewMessage(subject string, to string, tmplName string, data interface{}) (email.Message, error) {
	var buf bytes.Buffer
	err := e.templateCache.ExecuteTemplate(&buf, tmplName, data)
	if err != nil {
		return email.Message{}, err
	}
	msg := email.Message{
		Subject:     subject,
		To:          to,
		HtmlContent: buf.String(),
	}
	return msg, nil
}


func (e *Emailer) Send(ctx context.Context, m email.Message) error {
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
