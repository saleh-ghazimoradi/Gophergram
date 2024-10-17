package service

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"time"
)

const (
	FromName            = "Gophergram"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "template"
var FS embed.FS

type Mailer interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

type sendGridMailerService struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func (s *sendGridMailerService) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)
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
	for i := 0; i < maxRetries; i++ {
		response, retryErr := s.client.Send(message)
		if retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return response.StatusCode, nil
	}
	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetries, retryErr)
}

func NewSendGridMailer(apiKey, fromEmail string) Mailer {
	client := sendgrid.NewSendClient(apiKey)
	return &sendGridMailerService{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}
