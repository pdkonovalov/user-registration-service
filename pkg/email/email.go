package email

import (
	"net/smtp"

	"github.com/pdkonovalov/user-registration-service/pkg/config"
	"github.com/pdkonovalov/user-registration-service/pkg/email/templates"
)

type Email struct {
	email string
	host  string
	auth  smtp.Auth
}

func Init(config *config.Config) (*Email, error) {
	email := &Email{
		email: config.EmailAddres,
		host:  config.EmailHost,
		auth:  smtp.PlainAuth("", config.EmailAddres, config.EmailPassword, config.EmailHost),
	}
	err := email.Send(email.email, templates.StartServiceMsg(email.email))
	if err != nil {
		return nil, err
	}
	return email, nil
}

func (email *Email) Send(to string, msg string) error {
	return smtp.SendMail(email.host+":587", email.auth, email.email, []string{to}, []byte(msg))
}
