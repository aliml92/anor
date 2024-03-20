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

func (e *brevoEmailer) SendVerificationMessageWithOTP(ctx context.Context, otp string, email string) error {
	var tplBuffer bytes.Buffer
	tmpl := e.templateCache.Lookup(e.config.SignupVerificationTemplateName)

	err := tmpl.Execute(&tplBuffer, otp)
	if err != nil {
		return err
	}
	_, _, err = e.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Anor",
			Email: e.config.FromEmail,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: email,
			},
		},
		Subject:     "Anor account details confirmation",
		HtmlContent: tplBuffer.String(),
	})

	return err
}
