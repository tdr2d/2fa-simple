package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type MailUser struct {
	Name string
	Mail string
}

func SendGrid(from MailUser, to MailUser, subject string, content string) error {
	f := mail.NewEmail(from.Name, from.Mail)
	t := mail.NewEmail(to.Name, to.Mail)
	message := mail.NewSingleEmail(f, subject, t, "", content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if response != nil && response.StatusCode >= 400 {
		logrus.Info(message)
		logrus.Error(fmt.Sprintf("%d %s", response.StatusCode, response.Body))
	}

	return err
}
