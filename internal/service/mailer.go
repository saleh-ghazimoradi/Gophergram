package service

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"time"
)

//go:embed "template"
var FS embed.FS

type Mailer interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

type mailService struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func (m *mailService) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(config.AppConfig.Mail.FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "template/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < int(config.AppConfig.Mail.MaxRetries); i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return response.StatusCode, nil
	}
	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", int(config.AppConfig.Mail.MaxRetries), retryErr)
}

func NewMailer(fromEmail, apiKey string) Mailer {
	client := sendgrid.NewSendClient(apiKey)
	return &mailService{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}
