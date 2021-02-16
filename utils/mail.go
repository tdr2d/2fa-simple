package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// type Config struct {
// 	Host     string
// 	Port     string
// 	Password string
// }

// func SendHtmlMail(from string, to []string, subject string, html string, config Config) error {
// 	auth := smtp.PlainAuth("", from, config.Password, config.Host)

// 	var body bytes.Buffer
// 	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
// 	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

// 	return smtp.SendMail(config.Host+":"+config.Port, auth, from, to, []byte(html))
// }

type Recipient struct {
	Name string
	Mail string
}

func SendGrid(from Recipient, to Recipient, subject string, content string) error {
	f := mail.NewEmail(from.Name, from.Mail)
	t := mail.NewEmail(to.Name, to.Mail)
	message := mail.NewSingleEmail(f, subject, t, "", content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if response != nil {
		logrus.Info(fmt.Sprintf("%d %s", response.StatusCode, response.Body))
	}

	return err
}
